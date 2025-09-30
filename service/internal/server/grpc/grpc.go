// internal/server/grpc/grpc.go
package server_grpc

import (
	"service/internal/conf/v1"
	"service/internal/server/utils/requestlog"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// GRPCRegistrar is a function that registers routes on the server.
type GRPCRegister func(*grpc.Server)

func NewGRPCServer(
	c *conf.Server,
	registrers []GRPCRegister,
	// authGroups []endpoint.ServiceGroup,
	logger log.Logger,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
		// Add logging for unary and stream requests
		grpc.UnaryInterceptor(requestlog.UnaryLogInterceptor()),
		grpc.StreamInterceptor(requestlog.StreamLogInterceptor()),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	// Automatic registration of all modules
	for _, r := range registrers {
		r(srv)
	}

	return srv
}
