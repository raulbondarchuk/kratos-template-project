package requestlog

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	mylog "service/pkg/logger"
)

// getClientIP извлекает IP адрес клиента из разных типов запросов
func getClientIP(ctx context.Context, r *http.Request) string {
	// Для HTTP запросов
	if r != nil {
		// X-Forwarded-For: берем первый IP
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

	// Для gRPC запросов
	if p, ok := peer.FromContext(ctx); ok {
		if p.Addr != nil {
			return p.Addr.String()
		}
	}
	return "unknown"
}

// logRequest логирует информацию о запросе
func logRequest(method, path string, fields map[string]interface{}) {
	mylog.Route(method, path, fields)
}

// HTTP middleware

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

// HTTPLogMiddleware создает HTTP middleware для логирования запросов
func HTTPLogMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w}
			start := time.Now()

			next.ServeHTTP(sw, r)

			logRequest(
				r.Method,
				r.URL.Path,
				map[string]interface{}{
					"ip":      getClientIP(r.Context(), r),
					"status":  sw.status,
					"size":    sw.size,
					"latency": time.Since(start).String(),
				},
			)
		})
	}
}

// gRPC interceptors

// UnaryLogInterceptor создает unary interceptor для логирования gRPC запросов
func UnaryLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		// Выполняем запрос
		resp, err = handler(ctx, req)

		// Определяем статус
		status := "OK"
		if err != nil {
			status = err.Error()
		}

		logRequest(
			"gRPC",
			info.FullMethod,
			map[string]interface{}{
				"ip":      getClientIP(ctx, nil),
				"status":  status,
				"latency": time.Since(start).String(),
			},
		)

		return resp, err
	}
}

// StreamLogInterceptor создает stream interceptor для логирования gRPC стримов
func StreamLogInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Выполняем обработку стрима
		err := handler(srv, ss)

		// Определяем статус
		status := "OK"
		if err != nil {
			status = err.Error()
		}

		logRequest(
			"gRPC Stream",
			info.FullMethod,
			map[string]interface{}{
				"ip":      getClientIP(ss.Context(), nil),
				"status":  status,
				"latency": time.Since(start).String(),
			},
		)

		return err
	}
}
