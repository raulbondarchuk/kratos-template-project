package example

import (
	api_example "service/api/example/v1"
	example_service "service/internal/feature/example/v1/service"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Module-local types to avoid wire type collisions
type HTTPRegister func(*http.Server)
type GRPCRegister func(*grpc.Server)

// NOTE: versioned service interfaces from proto: Examplev1Service...
var _ api_example.Examplev1ServiceHTTPServer = (*example_service.ExampleService)(nil)
var _ api_example.Examplev1ServiceServer = (*example_service.ExampleService)(nil)

func NewExampleHTTPRegistrer(s api_example.Examplev1ServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api_example.RegisterExamplev1ServiceHTTPServer(srv, s)
	}
}

func NewExampleGRPCRegistrer(s api_example.Examplev1ServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api_example.RegisterExamplev1ServiceServer(srv, s)
	}
}
