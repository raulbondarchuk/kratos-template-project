// cmd/service/registrars_agg.go
package main

import (
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"
)

type AllRegistrers struct {
	HTTP []server_http.HTTPRegister
	GRPC []server_grpc.GRPCRegister
}

func BuildAllRegistrars(
	// HTTP
	templateHTTP server_http.HTTPRegister,
	// add other HTTP-registrers for modules here:
	// usersHTTP server_http.HTTPRegister,
	// alertsHTTP server_http.HTTPRegister,

	// gRPC
	templateGRPC server_grpc.GRPCRegister,
	// add other gRPC-registrers for modules here:
	// usersGRPC server_grpc.GRPCRegister,
	// alertsGRPC server_grpc.GRPCRegister,
) AllRegistrers {
	return AllRegistrers{
		HTTP: []server_http.HTTPRegister{
			templateHTTP,
			// usersHTTP, alertsHTTP, ...
		},
		GRPC: []server_grpc.GRPCRegister{
			templateGRPC,
			// usersGRPC, alertsGRPC, ...
		},
	}
}

// Two small providers for Wire:
func ProvideHTTPRegistrers(all AllRegistrers) []server_http.HTTPRegister { return all.HTTP }
func ProvideGRPCRegistrers(all AllRegistrers) []server_grpc.GRPCRegister { return all.GRPC }
