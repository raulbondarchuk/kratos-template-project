// cmd/service/registers_agg.go
package main

import (
	templatev1 "service/internal/feature/template/v1" // without /v — there is the root package template
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

)

type AllRegistrers struct {
	HTTP []server_http.HTTPRegister
	GRPC []server_grpc.GRPCRegister
}

func BuildAllRegistrars(
	// HTTP
	templatev1HTTP templatev1.HTTPRegister,
	// add other HTTP-registrers for modules here:

	// gRPC
	templatev1GRPC templatev1.GRPCRegister,
	// add other gRPC-registrers for modules here:

) AllRegistrers {
	return AllRegistrers{
		HTTP: []server_http.HTTPRegister{
			server_http.HTTPRegister(templatev1HTTP),
			// add other HTTP-registrers for modules here:

		},
		GRPC: []server_grpc.GRPCRegister{
			server_grpc.GRPCRegister(templatev1GRPC),
			// add other gRPC-registrers for modules here:

		},
	}
}

func ProvideHTTPRegistrers(all AllRegistrers) []server_http.HTTPRegister { return all.HTTP }
func ProvideGRPCRegistrers(all AllRegistrers) []server_grpc.GRPCRegister { return all.GRPC }

