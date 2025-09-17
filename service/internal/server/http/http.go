// internal/server/http/http.go
package server_http

import (
	"service/internal/conf/v1"
	"service/internal/middleware/requestlog"
	"service/internal/server/http/openapi/scalar"
	"service/internal/server/http/openapi/swagger"
	"service/internal/server/http/sys"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(
	c *conf.Server,
	registrers []HTTPRegister,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(recovery.Recovery()),
		http.Filter(requestlog.RequestLogFilter()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	// Automatic registration of all modules
	for _, r := range registrers {
		r(srv)
	}

	// Documentation and system endpoints
	swagger.AttachEmbeddedSwaggerUI(srv)
	scalar.AttachScalarDocs(srv)
	sys.LoadSystemEndpoints(srv)

	return srv
}
