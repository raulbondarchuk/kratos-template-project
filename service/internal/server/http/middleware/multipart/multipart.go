package multipart

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
)

/*
// Get files from multipart context
	m, ok := multipart.FromContext(ctx)
	if !ok || len(m.Files) == 0 {
		return &pb.UpsertExcelResponse{
			ProcessedRows: 0,
			Meta: &pb.MetaResponse{
				Code:    pb.ResponseCode_RESPONSE_CODE_BAD_REQUEST,
				Message: "No files were uploaded",
			},
		}, nil
	}

	// Get files with key "file"
	files, ok := m.Files["file"]
	if !ok || len(files) == 0 {
		return &pb.UpsertExcelResponse{
			ProcessedRows: 0,
			Meta: &pb.MetaResponse{
				Code:    pb.ResponseCode_RESPONSE_CODE_BAD_REQUEST,
				Message: "Field 'file' is required",
			},
		}, nil
	}
*/

type FileData struct {
	File    multipart.File
	Header  *multipart.FileHeader
	Content []byte
}

type Multipart struct {
	Files map[string][]*FileData // key - form field name, value - array of files
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

				m := &Multipart{
					Files: make(map[string][]*FileData),
				}

				// Process all files from the form
				for field, headers := range r.MultipartForm.File {
					files := make([]*FileData, 0, len(headers))

					for _, header := range headers {
						file, err := header.Open()
						if err != nil {
							http.Error(w, fmt.Sprintf("failed to open file %s: %v", header.Filename, err), http.StatusBadRequest)
							return
						}

						content, err := io.ReadAll(file)
						if err != nil {
							file.Close()
							http.Error(w, fmt.Sprintf("failed to read file %s: %v", header.Filename, err), http.StatusInternalServerError)
							return
						}
						file.Close() // Close immediately after reading

						files = append(files, &FileData{
							File:    file,
							Header:  header,
							Content: content,
						})
					}

					if len(files) > 0 {
						m.Files[field] = files
					}
				}

				// Check that at least one file was uploaded
				if len(m.Files) == 0 {
					http.Error(w, "no files were uploaded", http.StatusBadRequest)
					return
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
