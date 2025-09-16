// internal/server/grpc/registrar.go
package server_grpc

import "github.com/go-kratos/kratos/v2/transport/grpc"

// GRPCRegistrar is a function that registers routes on the server.
type GRPCRegister func(*grpc.Server)
