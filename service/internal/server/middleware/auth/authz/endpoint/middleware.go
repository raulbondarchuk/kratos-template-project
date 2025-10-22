package endpoint

import (
	"context"
	"fmt"

	http_errors "service/internal/server/http/middleware/errors"
	"service/internal/server/middleware/auth/auth/paseto"
	"service/pkg/logger"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// RoleMiddleware checks token and ensures user has at least one of requiredRoles.
// Enforced ONLY for HTTP; for gRPC it's skipped.
func RoleMiddleware(requiredRoles []string) middleware.Middleware {
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

			// HTTP: read Authorization
			token, err := GetAccessToken(ctx)
			if err != nil {
				return nil, http_errors.Unauthorized(ReasonAuthz, err.Error(), nil)
			}
			if token == "" {
				return nil, http_errors.Unauthorized(ReasonAuthz, ErrMissingAuthorizationHeader.Error(), nil)
			}

			// verify token
			claims, err := paseto.VerifyAccessToken(token)
			if err != nil {
				logger.Warn("RoleMiddleware: token verification failed",
					map[string]interface{}{"error": err})
				return nil, http_errors.Unauthorized(ReasonAuthz, fmt.Sprintf("invalid token: %v", err), nil)
			}

			// roles are CSV in claims.Roles
			userRoles := splitCSV(claims.Roles)
			if !HasRequiredRole(userRoles, requiredRoles) {
				logger.Warn("RoleMiddleware: insufficient permissions",
					map[string]interface{}{"required": requiredRoles, "got": claims.Roles})
				return nil, http_errors.Forbidden(
					ReasonAuthz,
					fmt.Sprintf("insufficient permissions: required one of %v, got %v", requiredRoles, claims.Roles),
					nil,
				)
			}

			// put roles/claims into ctx if needed
			ctx = context.WithValue(ctx, ctxKeyRoles, userRoles)
			ctx = context.WithValue(ctx, ctxKeyClaims, claims)

			return next(ctx, req)
		}
	}
}
