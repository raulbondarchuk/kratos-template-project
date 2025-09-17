// cmd/service/registers_agg.go
package main

import (
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"
)

// AllRegistrers Es la estructura que contiene los registros de los servidores HTTP y GRPC
type AllRegistrers struct {
	HTTP []server_http.HTTPRegister
	GRPC []server_grpc.GRPCRegister
}

func BuildAllRegistrars(
	// HTTP
	templateHTTP server_http.HTTPRegister,
	// add other HTTP-registrers for modules here:

	// gRPC
	templateGRPC server_grpc.GRPCRegister,
	// add other gRPC-registrers for modules here:

) AllRegistrers {
	return AllRegistrers{
		HTTP: []server_http.HTTPRegister{
			templateHTTP,
			// add other HTTP-registrers for modules here:

		},
		GRPC: []server_grpc.GRPCRegister{
			templateGRPC,
			// add other gRPC-registrers for modules here:

		},
	}
}

// Two small providers for Wire:
func ProvideHTTPRegistrers(all AllRegistrers) []server_http.HTTPRegister { return all.HTTP }
func ProvideGRPCRegistrers(all AllRegistrers) []server_grpc.GRPCRegister { return all.GRPC }
