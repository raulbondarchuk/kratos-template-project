package template

import (
	api_template "service/api/template/v1"
	template_service "service/internal/feature/template/v1/service"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Module types (other packages will have their own, so types are different for wire)
type HTTPRegister func(*http.Server)
type GRPCRegister func(*grpc.Server)

var _ api_template.Templatev1ServiceHTTPServer = (*template_service.TemplateService)(nil)
var _ api_template.Templatev1ServiceServer = (*template_service.TemplateService)(nil)

func NewTemplatesHTTPRegistrer(s api_template.Templatev1ServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api_template.RegisterTemplatev1ServiceHTTPServer(srv, s)
	}
}

func NewTemplatesGRPCRegistrer(s api_template.Templatev1ServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api_template.RegisterTemplatev1ServiceServer(srv, s)
	}
}
