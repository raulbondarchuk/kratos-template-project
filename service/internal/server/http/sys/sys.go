package sys

import (
	"encoding/json"
	stdhttp "net/http"
	"time"

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
