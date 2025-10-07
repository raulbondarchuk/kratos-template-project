package swagger

import (
	stdhttp "net/http"
	"strings"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func setOnlyContentType(w stdhttp.ResponseWriter, ct string) {
	h := w.Header()
	for k := range h {
		h.Del(k)
	}
	h.Set("Content-Type", ct)
}

func setNoCache(w stdhttp.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
}

func httpNotFound(w stdhttp.ResponseWriter) {
	setOnlyContentType(w, "text/plain; charset=utf-8")
	w.WriteHeader(stdhttp.StatusNotFound)
	_, _ = w.Write([]byte("not found"))
}

// "" or "/" → ""
func cleanBase(b string) string {
	b = strings.TrimSpace(b)
	if b == "" || b == "/" {
		return ""
	}
	if !strings.HasPrefix(b, "/") {
		b = "/" + b
	}
	b = strings.TrimRight(b, "/")
	return b
}

// Dynamic base from X-Forwarded-Prefix + cfg.Base
func baseForReq(r *stdhttp.Request, cfg *Config) string {
	if xf := strings.TrimSpace(r.Header.Get("X-Forwarded-Prefix")); xf != "" {
		return cleanBase(xf)
	}
	return cfg.Base
}

// Where to hang cookie within the current request
func cookiePathForReq(r *stdhttp.Request, cfg *Config) string {
	b := baseForReq(r, cfg)
	if b == "" {
		return "/docs/"
	}
	return b + "/docs/"
}

// Double registration of paths (with cfg.Base and without it)
func reg(s *kratoshttp.Server, base string, path string, h stdhttp.HandlerFunc) {
	if base != "" {
		s.HandleFunc(base+path, h)
	}
	s.HandleFunc(path, h)
}

func regFS(s *kratoshttp.Server, base string, path string, h stdhttp.Handler) {
	if base != "" {
		s.Handle(base+path, h)
	}
	s.Handle(path, h)
}
