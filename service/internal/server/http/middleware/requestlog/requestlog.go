// internal/middleware/requestlog/requestlog.go
package requestlog

import (
	"net"
	"net/http"
	"strings"
	"time"

	mylog "service/pkg/logger"
)

// wrapper over ResponseWriter to capture status and size
type statusWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func clientIP(r *http.Request) string {
	// X-Forwarded-For: take first IP
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	// Nginx/Traefik
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	// RemoteAddr -> host
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

// RequestLogFilter — net/http filter for Kratos HTTP server.
func RequestLogFilter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w}
			start := time.Now()

			next.ServeHTTP(sw, r)

			mylog.Route(
				r.Method,
				r.URL.Path,
				map[string]interface{}{
					"ip":      clientIP(r),
					"status":  sw.status,
					"size":    sw.size,
					"latency": time.Since(start).String(),
				},
			)
		})
	}
}
