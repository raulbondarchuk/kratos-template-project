// internal/server/middleware/traffic/mw_grpc.go
package traffic

import (
	"context"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

/*
   Middleware gRPC:
   Misma lógica que HTTP, pero claves derivadas de metadata/peer.
*/

func (b *Builder) GRPC() middleware.Middleware {
	var tail ratelimit.Limiter
	if b.cfg.EnableCPU {
		tail = bbr.NewLimiter(
			bbr.WithWindow(b.cfg.CPUWindow),
			bbr.WithBucket(b.cfg.CPUBuckets),
			bbr.WithCPUThreshold(b.cfg.CPUThreshold),
			bbr.WithCPUQuota(b.cfg.CPUQuota),
		)
	}

	b.cfg.LogHelper.Infof("[gRPC] [TRAFFIC RATE LIMIT] middleware initialized")

	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 1) InFlight
			leave, blocked := b.tryInflight()
			if blocked {
				return nil, tooMany()
			}
			defer leave()

			// 2) Clave (global/ip/user)
			md, _ := metadata.FromIncomingContext(ctx)
			var p *peer.Peer
			if pp, ok := peer.FromContext(ctx); ok {
				p = pp
			}
			key := "global"
			switch b.cfg.KeyBy {
			case KeyIP:
				key = "ip:" + ipFromGRPC(md, p)
			case KeyUser:
				uid := ""
				if vals := md.Get("x-user-id"); len(vals) > 0 {
					uid = vals[0]
				}
				if uid != "" {
					key = "u:" + uid
				} else {
					key = "ip:" + ipFromGRPC(md, p)
				}
			default:
				key = "global"
			}

			// 3) Rate
			if b.rl != nil && !b.rl.Allow(key) {
				return nil, tooMany()
			}

			// 4) BBR
			if tail != nil {
				if done, err := tail.Allow(); err != nil {
					return nil, tooMany()
				} else {
					defer done(ratelimit.DoneInfo{})
				}
			}

			return next(ctx, req)
		}
	}
}
