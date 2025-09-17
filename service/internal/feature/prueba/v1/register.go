package prueba

import (
	api_prueba "service/api/prueba/v1"
	prueba_service "service/internal/feature/prueba/v1/service"
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var _ api_prueba.PruebaServiceHTTPServer = (*prueba_service.PruebaService)(nil)
var _ api_prueba.PruebaServiceServer     = (*prueba_service.PruebaService)(nil)

func NewPruebaHTTPRegister(s api_prueba.PruebaServiceHTTPServer) server_http.HTTPRegister {
	return func(srv *http.Server) {
		api_prueba.RegisterPruebaServiceHTTPServer(srv, s)
	}
}

func NewPruebaGRPCRegister(s api_prueba.PruebaServiceServer) server_grpc.GRPCRegister {
	return func(srv *grpc.Server) {
		api_prueba.RegisterPruebaServiceServer(srv, s)
	}
}