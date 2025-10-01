// cmd/service/registers_agg.go
package main

import (

	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	examplev1 "service/internal/feature/example/v1"

)

type AllRegistrers struct {
	HTTP []server_http.HTTPRegister
	GRPC []server_grpc.GRPCRegister
}

func BuildAllRegistrars(
	// HTTP

	// add other HTTP-registrers for modules here:

	examplev1HTTP examplev1.HTTPRegister,

	// gRPC

	// add other gRPC-registrers for modules here:

	examplev1GRPC examplev1.GRPCRegister,

) AllRegistrers {
	return AllRegistrers{
		HTTP: []server_http.HTTPRegister{

			// add other HTTP-registrers for modules here:

			server_http.HTTPRegister(examplev1HTTP),

		},
		GRPC: []server_grpc.GRPCRegister{

			// add other gRPC-registrers for modules here:

			server_grpc.GRPCRegister(examplev1GRPC),

		},
	}
}

func ProvideHTTPRegistrers(all AllRegistrers) []server_http.HTTPRegister { return all.HTTP }
func ProvideGRPCRegistrers(all AllRegistrers) []server_grpc.GRPCRegister { return all.GRPC }

