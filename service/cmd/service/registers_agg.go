// cmd/service/registers_agg.go
package main

import (
	examplev1 "service/internal/feature/example/v1" // without /v — there is the root package example
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"
)

type AllRegistrers struct {
	HTTP []server_http.HTTPRegister
	GRPC []server_grpc.GRPCRegister
}

func BuildAllRegistrars(
	// HTTP
	examplev1HTTP examplev1.HTTPRegister,
	// add other HTTP-registrers for modules here:

	// gRPC
	examplev1GRPC examplev1.GRPCRegister,
	// add other gRPC-registrers for modules here:

) AllRegistrers {
	return AllRegistrers{
		HTTP: []server_http.HTTPRegister{
			server_http.HTTPRegister(examplev1HTTP),
			// add other HTTP-registrers for modules here:

		},
		GRPC: []server_grpc.GRPCRegister{
			server_grpc.GRPCRegister(examplev1GRPC),
			// add other gRPC-registrers for modules here:

		},
	}
}

func ProvideHTTPRegistrers(all AllRegistrers) []server_http.HTTPRegister { return all.HTTP }
func ProvideGRPCRegistrers(all AllRegistrers) []server_grpc.GRPCRegister { return all.GRPC }
