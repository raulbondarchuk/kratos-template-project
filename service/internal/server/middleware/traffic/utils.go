package traffic

import (
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// helpers for key extraction

// HTTP: real client IP (proxies aware), fallback to RemoteAddr
func ipFromRequest(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host == "" {
		return "unknown"
	}
	return host
}

// HTTP: user id from header (customize if you have signed claims)
func userFromRequest(r *http.Request) string {
	if uid := r.Header.Get("X-User-ID"); uid != "" {
		return uid
	}
	return ""
}

// gRPC: client IP from metadata/peer
func ipFromGRPC(md metadata.MD, p *peer.Peer) string {
	if vals := md.Get("x-real-ip"); len(vals) > 0 && vals[0] != "" {
		return vals[0]
	}
	if vals := md.Get("x-forwarded-for"); len(vals) > 0 && vals[0] != "" {
		xff := vals[0]
		if i := strings.IndexByte(xff, ','); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if p != nil && p.Addr != nil {
		host, _, _ := net.SplitHostPort(p.Addr.String())
		return host
	}
	return "unknown"
}
