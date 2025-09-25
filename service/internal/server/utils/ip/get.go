package server_utils_ip

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func GetIP(ctx context.Context) string {
	ip := "unknown"
	if tr, ok := transport.FromServerContext(ctx); ok {
		switch tr.Kind() {
		case transport.KindHTTP:
			if ht, ok := tr.(*khttp.Transport); ok {
				ip = extractIP(ht.Request())
			}
		case transport.KindGRPC:
			// For gRPC requests, we get the IP from the headers
			ip = tr.RequestHeader().Get("x-real-ip")
			if ip == "" {
				ip = tr.RequestHeader().Get("x-forwarded-for")
			}
		}
	}
	return ip
}

// extractIP extracts the IP address from the HTTP request with consideration for proxies
func extractIP(r *http.Request) string {
	// X-Real-IP
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// X-Forwarded-For
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// We get the first IP from the list (the most original client IP)
		if i := strings.Index(ip, ","); i > 0 {
			ip = ip[:i]
		}
		return ip
	}

	// RemoteAddr
	ip = r.RemoteAddr
	if ip != "" {
		// We remove the port if it exists
		if i := strings.LastIndex(ip, ":"); i > 0 {
			ip = ip[:i]
		}
		return ip
	}

	return "unknown"
}
