package authz

import (
	"service/internal/server/middleware/auth/authz/endpoint"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ProviderSet creates a auth middleware for HTTP server.
func ProviderSet(groups []endpoint.ServiceGroup) middleware.Middleware {
	return endpoint.CreateMiddleware(groups)
}
