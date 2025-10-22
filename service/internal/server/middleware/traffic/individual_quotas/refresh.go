package individual_quotas

import (
	"context"

	"golang.org/x/time/rate"
)

/*
   Refresco de cuotas: pedir al servicio externo, validar respuesta,
   reconstruir la tabla y limiters, o vaciar todo en caso de error.
*/

// refreshOnce baja cuotas y actualiza limiters.
// Si falla, se deja la tabla vacía (y se limpian limiters).
func (iq *IQ) RefreshOnce(ctx context.Context) {
	if iq.project == "" {
		// nada que hacer si no hay proyecto
		return
	}

	var body quotaResponse
	resp, err := iq.rest.R().
		SetContext(ctx).
		SetQueryParam("project", iq.project).
		SetResult(&body).
		Get(iq.serviceURL + "/qt")

	// Error de red / contexto / decode o HTTP no-2xx: vaciamos cuotas
	if err != nil {
		iq.logHelper.Errorf("[%s] [IQ] fetch error: %v", iq.serverType, err)
		iq.setEmptyQuotas()
		return
	}
	if resp == nil || resp.IsError() {
		status := ""
		if resp != nil {
			status = resp.Status()
		}
		iq.logHelper.Errorf("[%s] [IQ] bad HTTP status: %s", iq.serverType, status)
		iq.setEmptyQuotas()
		return
	}

	// Construimos la nueva tabla
	next := make(map[string]quotaCfg, len(body.Items))
	for _, it := range body.Items {
		route := normRoute(it.Route)
		if route == "" || it.Quota <= 0 {
			continue
		}
		intSec := normIntervalSec(it.Interval)
		next[route] = quotaCfg{Quota: it.Quota, Interval: intSec}
		iq.ensureLimiter(route, it.Quota, intSec)
	}

	// Limpiar limiters que ya no existen
	iq.gcLimitersKeys(next)

	// Publicar la nueva tabla
	iq.quotas.Store(next)
	iq.logHelper.Infof("[%s] [IQ] quotas applied: %d routes", iq.serverType, len(next))
}

// setEmptyQuotas borra la tabla de cuotas y limpia todos los limiters.
func (iq *IQ) setEmptyQuotas() {
	iq.quotas.Store(make(map[string]quotaCfg))
	iq.mu.Lock()
	iq.limiters = make(map[string]*rate.Limiter)
	iq.mu.Unlock()
	iq.logHelper.Warnf("[%s] [IQ] quotas cleared: using empty set (fail-open)", iq.serverType)
}
