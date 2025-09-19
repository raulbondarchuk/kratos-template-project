package prueba

import (
	api_prueba "service/api/prueba/v1"
	prueba_service "service/internal/feature/prueba/v1/service"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Module-local types to avoid wire type collisions
type HTTPRegister func(*http.Server)
type GRPCRegister func(*grpc.Server)

// NOTE: versioned service interfaces from proto: Pruebav1Service...
var _ api_prueba.Pruebav1ServiceHTTPServer = (*prueba_service.PruebaService)(nil)
var _ api_prueba.Pruebav1ServiceServer     = (*prueba_service.PruebaService)(nil)

func NewPruebaHTTPRegistrer(s api_prueba.Pruebav1ServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api_prueba.RegisterPruebav1ServiceHTTPServer(srv, s)
	}
}

func NewPruebaGRPCRegistrer(s api_prueba.Pruebav1ServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api_prueba.RegisterPruebav1ServiceServer(srv, s)
	}
}