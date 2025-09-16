// internal/feature/template/registrars.go
package template

import (
	api_template "service/api/template/v1"
	template_service "service/internal/feature/template/service"
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var _ api_template.TemplatesServiceHTTPServer = (*template_service.TemplatesService)(nil)
var _ api_template.TemplatesServiceServer = (*template_service.TemplatesService)(nil)

// HTTP
func NewTemplatesHTTPRegistrer(s api_template.TemplatesServiceHTTPServer) server_http.HTTPRegister {
	return func(srv *http.Server) {
		api_template.RegisterTemplatesServiceHTTPServer(srv, s)
	}
}

func NewTemplatesGRPCRegistrer(s api_template.TemplatesServiceServer) server_grpc.GRPCRegister {
	return func(srv *grpc.Server) {
		api_template.RegisterTemplatesServiceServer(srv, s)
	}
}
