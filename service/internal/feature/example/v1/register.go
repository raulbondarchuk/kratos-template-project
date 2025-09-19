package example

import (
	api "service/api/example/v1"
	service "service/internal/feature/example/v1/service"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Module types (other packages will have their own, so types are different for wire)
type HTTPRegister func(*http.Server)
type GRPCRegister func(*grpc.Server)

var _ api.Examplev1ServiceHTTPServer = (*service.ExampleService)(nil)
var _ api.Examplev1ServiceServer = (*service.ExampleService)(nil)

func NewExampleHTTPRegistrer(s api.Examplev1ServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api.RegisterExamplev1ServiceHTTPServer(srv, s)
	}
}

func NewExampleGRPCRegistrer(s api.Examplev1ServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api.RegisterExamplev1ServiceServer(srv, s)
	}
}
