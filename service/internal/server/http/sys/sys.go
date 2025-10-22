package sys

import (
	"encoding/json"
	"net/http"
	stdhttp "net/http"
	"time"

	iqpkg "service/internal/server/middleware/traffic/individual_quotas"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime = time.Now()

func LoadSystemEndpoints(srv *khttp.Server) {
	// Prometheus metrics
	srv.Handle("/metrics", promhttp.Handler())

	// Simple health
	srv.HandleFunc("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data := map[string]any{
			"time":   time.Now().UTC().Format(time.RFC3339),
			"uptime": time.Since(startTime).String(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(data)
	})

}

func LoadQuotasRefreshEndpoint(srv *khttp.Server, iq *iqpkg.IQ) {
	if iq == nil {
		return
	}

	srv.HandleFunc("/qt", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("method not allowed"))
			return
		}

		start := time.Now()
		iq.RefreshOnce(r.Context())
		elapsed := time.Since(start)

		resp := map[string]any{
			"ok":            true,
			"refreshed_at":  time.Now().UTC().Format(time.RFC3339),
			"routes_loaded": iq.QuotasLen(),
			"took_ms":       elapsed.Milliseconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}
