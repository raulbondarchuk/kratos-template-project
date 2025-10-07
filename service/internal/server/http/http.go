// internal/server/http/http.go
package server_http

import (
	openapifs "service/docs"
	"service/internal/conf/v1"
	"service/internal/server/http/middleware/multipart"
	"service/internal/server/http/openapi/swagger"
	"service/internal/server/http/sys"
	"service/internal/server/middleware/auth/authz"
	"service/internal/server/middleware/auth/authz/endpoint"
	"service/internal/server/utils/requestlog"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// HTTPRegistrar is a function that registers routes on the server.
type HTTPRegister func(*http.Server)

func NewHTTPServer(
	c *conf.Server,
	app *conf.App,
	registrers []HTTPRegister,
	authGroups []endpoint.ServiceGroup,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			authz.ProviderSet(authGroups), // add auth middleware for roles
			multipart.Middleware(32<<20),  // 32MB max memory for file uploads
		),
		http.Filter(requestlog.HTTPLogMiddleware()),
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

	// 1) http
	swagger.AttachEmbeddedSwaggerUIWithConfig(srv, swagger.Config{
		Base:          "",
		DocsFS:        openapifs.FS,
		CookieName:    "swagger_default",
		ProjectPrefix: "/" + app.GetName(),
		ServiceName:   app.GetName(),
	})

	// 2) https
	swagger.AttachEmbeddedSwaggerUIWithConfig(srv, swagger.Config{
		Base:          "/" + app.GetName(),
		DocsFS:        openapifs.FS,
		CookieName:    "swagger_" + app.GetName(),
		ProjectPrefix: "/" + app.GetName(),
		ServiceName:   app.GetName(),
	})

	sys.LoadSystemEndpoints(srv)

	return srv
}
