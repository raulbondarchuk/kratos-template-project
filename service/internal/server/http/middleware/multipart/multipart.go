package multipart

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type FileData struct {
	File    multipart.File
	Header  *multipart.FileHeader
	Content []byte
}
type Multipart struct {
	Files map[string][]*FileData
}
type ctxKey struct{}

func FromContext(ctx context.Context) (*Multipart, bool) {
	m, ok := ctx.Value(ctxKey{}).(*Multipart)
	return m, ok
}

func Middleware(maxMemory int64) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, _ := transport.FromServerContext(ctx)
			htr, ok := tr.(*khttp.Transport)
			if !ok || htr == nil || htr.Request() == nil {
				return next(ctx, req)
			}

			r := htr.Request()
			if !isMultipartRequest(r) {
				return next(ctx, req)
			}

			if err := r.ParseMultipartForm(maxMemory); err != nil {
				// Convierte en 400 utilizando mi error middleware
				return nil, fmt.Errorf("failed to parse multipart form: %w", err)
			}
			defer func() {
				if r.MultipartForm != nil {
					_ = r.MultipartForm.RemoveAll()
				}
			}()

			m := &Multipart{Files: map[string][]*FileData{}}
			for field, headers := range r.MultipartForm.File {
				files := make([]*FileData, 0, len(headers))
				for _, hdr := range headers {
					f, err := hdr.Open()
					if err != nil {
						return nil, fmt.Errorf("failed to open file %s: %w", hdr.Filename, err)
					}
					b, err := io.ReadAll(f)
					_ = f.Close()
					if err != nil {
						return nil, fmt.Errorf("failed to read file %s: %w", hdr.Filename, err)
					}
					files = append(files, &FileData{Header: hdr, Content: b})
				}
				if len(files) > 0 {
					m.Files[field] = files
				}
			}

			// IMPORTANT: put in Kratos-context, which goes further in handler
			ctx = context.WithValue(ctx, ctxKey{}, m)
			return next(ctx, req)
		}
	}
}

func isMultipartRequest(r *http.Request) bool {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		return false
	}
	ct := strings.ToLower(r.Header.Get("Content-Type"))
	return strings.HasPrefix(ct, "multipart/")
}
