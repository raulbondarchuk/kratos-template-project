// internal/server/middleware/traffic/mw_http.go
package traffic

import (
	"context"
	"net/http"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	kratosErr "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

/*
   Middleware HTTP:
   - InFlight -> 429 si no hay slots.
   - TokenBucket (RPS/Burst) por clave -> 429 si excede.
   - Quota por ventana -> 429 + Retry-After.
   - Cola adaptativa BBR (CPU) -> 429 cuando sistema está saturado.
*/

func (b *Builder) HTTP() middleware.Middleware {
	var tail ratelimit.Limiter
	if b.cfg.EnableCPU {
		tail = bbr.NewLimiter(
			bbr.WithWindow(b.cfg.CPUWindow),
			bbr.WithBucket(b.cfg.CPUBuckets),
			bbr.WithCPUThreshold(b.cfg.CPUThreshold),
			bbr.WithCPUQuota(b.cfg.CPUQuota),
		)
	}

	b.cfg.LogHelper.Infof("[HTTP] [TRAFFIC RATE LIMIT] middleware initialized")

	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			hreq, ok := khttp.RequestFromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}

			// 1) InFlight
			leave, blocked := b.tryInflight()
			if blocked {
				return nil, tooMany()
			}
			defer leave()

			// 2) Clave (global/ip/user)
			key := "global"
			switch b.cfg.KeyBy {
			case KeyIP:
				key = "ip:" + ipFromRequest(hreq)
			case KeyUser:
				if u := userFromRequest(hreq); u != "" {
					key = "u:" + u
				} else {
					key = "ip:" + ipFromRequest(hreq)
				}
			default:
				key = "global"
			}

			// 3) Rate
			if b.rl != nil && !b.rl.Allow(key) {
				return nil, tooMany()
			}

			// 4) Cola BBR (CPU)
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

func tooMany() error {
	return kratosErr.New(http.StatusTooManyRequests, "RATE_LIMITED", "too many requests")
}
