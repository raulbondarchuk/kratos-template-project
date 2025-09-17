package template

import (
	api_template "service/api/template/v1"
	template_service "service/internal/feature/template/service"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Module types (other packages will have their own, so types are different for wire)
type HTTPRegister func(*http.Server)
type GRPCRegister func(*grpc.Server)

var _ api_template.TemplatesServiceHTTPServer = (*template_service.TemplatesService)(nil)
var _ api_template.TemplatesServiceServer = (*template_service.TemplatesService)(nil)

func NewTemplatesHTTPRegistrer(s api_template.TemplatesServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api_template.RegisterTemplatesServiceHTTPServer(srv, s)
	}
}

func NewTemplatesGRPCRegistrer(s api_template.TemplatesServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api_template.RegisterTemplatesServiceServer(srv, s)
	}
}
