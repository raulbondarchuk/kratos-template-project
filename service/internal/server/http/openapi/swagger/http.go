package swagger

import (
	stdhttp "net/http"
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
