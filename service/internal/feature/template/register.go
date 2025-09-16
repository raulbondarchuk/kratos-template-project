// internal/feature/template/registrars.go
package template

import (
	api_template "service/api/template"
	template_service "service/internal/feature/template/service"
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var _ api_template.TemplatesHTTPServer = (*template_service.TemplatesService)(nil)
var _ api_template.TemplatesServer = (*template_service.TemplatesService)(nil)

// HTTP
func NewTemplatesHTTPRegistrer(s api_template.TemplatesHTTPServer) server_http.HTTPRegister {
	return func(srv *http.Server) {
		api_template.RegisterTemplatesHTTPServer(srv, s)
	}
}

func NewTemplatesGRPCRegistrer(s api_template.TemplatesServer) server_grpc.GRPCRegister {
	return func(srv *grpc.Server) {
		api_template.RegisterTemplatesServer(srv, s)
	}
}
