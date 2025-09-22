package multipart

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
)

type Multipart struct {
	File    multipart.File
	Header  *multipart.FileHeader
	Content []byte
}

type multipartKey struct{}

// FromContext returns the Multipart from ctx.
func FromContext(ctx context.Context) (*Multipart, bool) {
	m, ok := ctx.Value(multipartKey{}).(*Multipart)
	return m, ok
}

// NewContext returns a new Context that carries Multipart.
func NewContext(ctx context.Context, m *Multipart) context.Context {
	return context.WithValue(ctx, multipartKey{}, m)
}

// Server is a multipart middleware.
func Server(maxMemory int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isMultipartRequest(r) {
				if err := r.ParseMultipartForm(maxMemory); err != nil {
					http.Error(w, fmt.Sprintf("failed to parse multipart form: %v", err), http.StatusBadRequest)
					return
				}
				defer func() {
					if r.MultipartForm != nil {
						r.MultipartForm.RemoveAll()
					}
				}()

				file, header, err := r.FormFile("file")
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to get file: %v", err), http.StatusBadRequest)
					return
				}
				defer file.Close()

				content, err := io.ReadAll(file)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to read file: %v", err), http.StatusInternalServerError)
					return
				}

				m := &Multipart{
					File:    file,
					Header:  header,
					Content: content,
				}

				ctx := NewContext(r.Context(), m)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isMultipartRequest(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return r.Method == "POST" && len(contentType) > 0 && contentType[:9] == "multipart"
}

// TransportKey is multipart transport key.
type TransportKey struct{}

// WithTransport returns a new context with Transport.
func WithTransport(ctx context.Context, tr transport.Transporter) context.Context {
	return context.WithValue(ctx, TransportKey{}, tr)
}
