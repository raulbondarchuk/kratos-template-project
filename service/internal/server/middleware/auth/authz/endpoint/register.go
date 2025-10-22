package endpoint

import (
	"context"
	"fmt"
	"strings"

	"service/pkg/logger"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// CreateMiddleware builds a middleware that, per-method, applies RoleMiddleware.
// Enforced ONLY for HTTP; gRPC is not enforced.
func CreateMiddleware(groups []ServiceGroup) middleware.Middleware {
	// map: MethodName -> required roles
	methodRoles := make(map[string][]string)

	for _, group := range groups {
		for _, m := range group.Methods {
			name := GetMethodName(m.Method) // e.g. "GetCompany"
			if name == "" {
				continue
			}
			logger.Debug(fmt.Sprintf("Registering roles for method %s", name), map[string]interface{}{"roles": m.RequiredRoles})
			methodRoles[name] = m.RequiredRoles
		}
	}

	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}
			// Skip for gRPC completely
			if tr.Kind() == transport.KindGRPC {
				return next(ctx, req)
			}

			// HTTP: figure out method name from operation
			op := tr.Operation() // e.g. "/pkg.Service/Method" or similar
			parts := strings.Split(op, "/")
			methodName := parts[len(parts)-1] // "Method"

			logger.Debug("Checking method", map[string]interface{}{"operation": op, "method": methodName})

			roles, exists := methodRoles[methodName]
			if !exists {
				// fallback: maybe someone stored full op string as key
				roles, exists = methodRoles[op]
			}

			if exists {
				logger.Debug("Found roles for method", map[string]interface{}{
					"method": methodName,
					"roles":  roles,
				})
				return RoleMiddleware(roles)(next)(ctx, req)
			}

			logger.Debug("No roles found for method", map[string]interface{}{"method": methodName})
			return next(ctx, req)
		}
	}
}
