package individual_quotas

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

/*
   Núcleo de IQ: tabla atómica de cuotas (quota + interval),
   y almacén de token-buckets por ruta.

   El refresco periódico de cuotas está implementado en refresh.go
*/

type ServerType string

const (
	HTTP ServerType = "HTTP"
	GRPC ServerType = "gRPC"
)

type IQ struct {
	// Config efectiva (resuelta con defaults)
	serviceURL   string
	project      string
	refreshEvery time.Duration
	burstFactor  float64
	strictMatch  bool

	// HTTP client
	rest *resty.Client

	// quotas {ruta → (quota, interval)} (lectura atómica, escritura reemplazando)
	quotas atomicMap[string, quotaCfg]

	// almacén de limiters por ruta (con mutex simple)
	mu       sync.Mutex
	limiters map[string]*rate.Limiter

	stopCh chan struct{}

	serverType ServerType
	logHelper  *log.Helper
}

// Config de cuota por ruta: N peticiones en una ventana de 'Interval' segundos.
type quotaCfg struct {
	Quota    int // número de peticiones permitidas por ventana
	Interval int // ventana en segundos (>0)
}

// New crea IQ con valores por defecto; `project` es obligatorio.
func New(project string, serverType ServerType, logger log.Logger) *IQ {
	if project == "" {
		panic("[INDIVIDUAL_QUOTAS] project is required")
	}

	url := defaultServiceURL
	refresh := defaultRefreshEvery
	burst := defaultBurstFactor
	strict := defaultStrictMatch

	cli := resty.New()
	cli.SetRetryCount(1)

	iq := &IQ{
		serviceURL:   strings.TrimRight(url, "/"),
		project:      project,
		serverType:   serverType,
		refreshEvery: refresh,
		burstFactor:  burst,
		strictMatch:  strict,
		rest:         cli,
		limiters:     make(map[string]*rate.Limiter),
		stopCh:       make(chan struct{}),
	}
	iq.logHelper = log.NewHelper(logger)
	iq.quotas.Store(make(map[string]quotaCfg)) // mapa vacío inicial
	return iq
}

// Start arranca el refresco (inmediato + cada refreshEvery).
// La lógica de refresco está en refresh.go (método refreshOnce).
func (iq *IQ) Start(ctx context.Context) {
	go func() {
		iq.RefreshOnce(ctx) // primer intento al arrancar

		t := time.NewTicker(iq.refreshEvery)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-iq.stopCh:
				return
			case <-t.C:
				iq.RefreshOnce(ctx)
			}
		}
	}()
}

// Stop detiene el bucle de refresco.
func (iq *IQ) Stop() { close(iq.stopCh) }

// QuotasLen devuelve el número de rutas con cuota cargada actualmente.
func (iq *IQ) QuotasLen() int {
	qm := iq.quotas.Load()
	return len(qm)
}

// ensureLimiter crea/actualiza el limiter para la ruta dada.
// Interpreta la cuota como "Quota peticiones por Interval segundos".
func (iq *IQ) ensureLimiter(route string, quota, intervalSec int) {
	// velocidad (eventos/segundo) para el token-bucket
	ratePerSec := float64(quota) / float64(intervalSec)

	// burst: capacidad del bucket. Usamos 'quota * BurstFactor' para permitir picos cortos.
	desiredR := rate.Limit(ratePerSec)
	desiredB := int(maxInt(1, int(float64(quota)*iq.burstFactor)))

	iq.mu.Lock()
	defer iq.mu.Unlock()

	if lim, ok := iq.limiters[route]; ok {
		if lim.Limit() != desiredR || lim.Burst() != desiredB {
			iq.limiters[route] = rate.NewLimiter(desiredR, desiredB)
		}
		return
	}
	iq.limiters[route] = rate.NewLimiter(desiredR, desiredB)
}

// gcLimitersKeys elimina limiters de rutas que ya no están en la tabla.
func (iq *IQ) gcLimitersKeys(next map[string]quotaCfg) {
	iq.mu.Lock()
	defer iq.mu.Unlock()
	for r := range iq.limiters {
		if _, keep := next[r]; !keep {
			delete(iq.limiters, r)
		}
	}
}

// getLimiter obtiene el limiter para una ruta (solo lectura).
func (iq *IQ) getLimiter(route string) *rate.Limiter {
	iq.mu.Lock()
	defer iq.mu.Unlock()
	return iq.limiters[route]
}

// allow chequea si un path/método puede pasar según la cuota/limiter.
func (iq *IQ) allow(route string) bool {
	route = normRoute(route)
	if route == "" {
		return true
	}

	qm := iq.quotas.Load()
	if len(qm) == 0 {
		return true
	}

	// match exacto o por prefijo (mejor coincidencia)
	if iq.strictMatch {
		if _, ok := qm[route]; !ok {
			return true
		}
		lim := iq.getLimiter(route)
		if lim == nil {
			return true
		}
		return lim.Allow()
	}

	best := ""
	for k := range qm {
		if strings.HasPrefix(route, k) && len(k) > len(best) {
			best = k
		}
	}
	if best == "" {
		return true
	}
	lim := iq.getLimiter(best)
	if lim == nil {
		return true
	}
	return lim.Allow()
}

// helpers locales (pueden vivir también en util.go)

func normIntervalSec(v int) int {
	if v <= 0 {
		return 1
	}
	return v
}

// normRoute normaliza: sin query, sin slashes extra, con "/" inicial, sin trailing "/".
func normRoute(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	for len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	return p
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/* -------- Mapa atómico genérico -------- */

type atomicMap[K comparable, V any] struct {
	ptr atomic.Pointer[map[K]V]
}

func (a *atomicMap[K, V]) Load() map[K]V {
	m := a.ptr.Load()
	if m == nil {
		return nil
	}
	return *m
}

func (a *atomicMap[K, V]) Store(m map[K]V) {
	a.ptr.Store(&m)
}
