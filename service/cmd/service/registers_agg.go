// cmd/service/registers_agg.go
package main

import (
	template "service/internal/feature/template" // without /v — there is the root package template
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"
)

type AllRegistrers struct {
	HTTP []server_http.HTTPRegister
	GRPC []server_grpc.GRPCRegister
}

func BuildAllRegistrars(
	// HTTP
	templateHTTP template.HTTPRegister,
	// add other HTTP-registrers for modules here:

	// gRPC
	templateGRPC template.GRPCRegister,
	// add other gRPC-registrers for modules here:

) AllRegistrers {
	return AllRegistrers{
		HTTP: []server_http.HTTPRegister{
			server_http.HTTPRegister(templateHTTP),
			// add other HTTP-registrers for modules here:
		},
		GRPC: []server_grpc.GRPCRegister{
			server_grpc.GRPCRegister(templateGRPC),
			// add other gRPC-registrers for modules here:
		},
	}
}

func ProvideHTTPRegistrers(all AllRegistrers) []server_http.HTTPRegister { return all.HTTP }
func ProvideGRPCRegistrers(all AllRegistrers) []server_grpc.GRPCRegister { return all.GRPC }
