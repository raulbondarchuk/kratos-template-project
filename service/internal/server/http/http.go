package server_http

import (
	"service/internal/conf"
	template_service "service/internal/feature/template/service"
	"service/internal/middleware/requestlog"
	"service/internal/server/http/openapi/scalar"
	"service/internal/server/http/openapi/swagger"
	"service/internal/server/http/sys"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server. FABRIC
func NewHTTPServer(c *conf.Server,
	// API modules
	template *template_service.TemplatesService,

	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
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

	// RequestLogFilter
	opts = append(opts, http.Filter(requestlog.RequestLogFilter()))

	srv := http.NewServer(opts...)

	// =================================================
	// === Register services ===========================
	// =================================================

	LoadRoutes(srv, template)

	// =================================================
	// =================================================

	// OPEN API
	swagger.AttachEmbeddedSwaggerUI(srv) // http://localhost:<PORT>/swagger-ui
	scalar.AttachScalarDocs(srv)         // http://localhost:<PORT>/docs

	// system functions
	sys.LoadSystemEndpoints(srv)

	return srv
}
