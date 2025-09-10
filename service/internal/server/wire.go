package server

import (
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	server_grpc.NewGRPCServer,
	server_http.NewHTTPServer,
)
