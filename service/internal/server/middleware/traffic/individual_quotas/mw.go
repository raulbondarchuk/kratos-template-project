package individual_quotas

import (
	"context"
	"net/http"

	kratosErr "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/metadata"
)

/*
   Middleware HTTP: aplica cuota por endpoint (RPS).
   Si no hay proyecto configurado — no-op.
*/

func (iq *IQ) HTTP() middleware.Middleware {
	if iq.project == "" {
		return passthrough()
	}
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			hreq, ok := khttp.RequestFromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}
			path := hreq.URL.Path
			if !iq.allow(path) {
				return nil, kratosErr.New(http.StatusTooManyRequests, "IQ_RATE_LIMITED", "too many requests for this endpoint")
			}
			return next(ctx, req)
		}
	}
}

/*
   Middleware gRPC: limita por método full path (ej.: "/pkg.Service/Method").
   Si no hay proyecto configurado — no-op.
*/

func (iq *IQ) GRPC() middleware.Middleware {
	if iq.project == "" {
		return passthrough()
	}
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			method := grpcMethod(ctx)
			if method != "" && !iq.allow(method) {
				return nil, kratosErr.New(http.StatusTooManyRequests, "IQ_RATE_LIMITED", "too many requests for this endpoint")
			}
			return next(ctx, req)
		}
	}
}

// grpcMethod intenta extraer ":path" del metadata entrante (p.ej. "/pkg.Svc/Method").
func grpcMethod(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v := md.Get(":path"); len(v) > 0 && v[0] != "" {
			return v[0]
		}
	}
	return ""
}

func passthrough() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return next(ctx, req)
		}
	}
}
