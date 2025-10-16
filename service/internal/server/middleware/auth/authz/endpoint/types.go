package endpoint

import (
	"context"
	"reflect"
	"runtime"
	"service/internal/server/middleware/auth/auth/paseto"
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/metadata"
)

// ----- service descriptions -----

type ServiceMethod struct {
	Service       interface{} // service instance
	Method        interface{} // method (func)
	RequiredRoles []string    // required roles
}

type ServiceGroup struct {
	Name    string
	Methods []ServiceMethod
}

func NewServiceMethod(service interface{}, method interface{}, roles ...string) ServiceMethod {
	return ServiceMethod{
		Service:       service,
		Method:        method,
		RequiredRoles: roles,
	}
}

// ----- context helpers -----

type ctxKey string

const (
	ctxKeyRoles       ctxKey = "roles"
	ctxKeyClaims      ctxKey = "claims"
	ctxKeyAccessToken ctxKey = "access_token"

	ctxKeyCompanyID ctxKey = "company_id"
	ctxKeyCliUser   ctxKey = "cliuser"
)

// TokenFromContext returns raw token stored by middleware (without "Bearer ").
func AccessTokenFromContext(ctx context.Context) string {
	if v := ctx.Value(ctxKeyAccessToken); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RolesFromContext returns roles previously stored by middleware.
func RolesFromContext(ctx context.Context) []string {
	if v := ctx.Value(ctxKeyRoles); v != nil {
		if roles, ok := v.([]string); ok {
			return roles
		}
	}
	return nil
}

// ClaimsFromContext returns claims previously stored by middleware.
func ClaimsFromContext(ctx context.Context) *paseto.Claims {
	if v := ctx.Value(ctxKeyClaims); v != nil {
		if claims, ok := v.(*paseto.Claims); ok {
			return claims
		}
	}
	return nil
}

func CompanyIDFromContext(ctx context.Context) uint {
	if v := ctx.Value(ctxKeyCompanyID); v != nil {
		if id, ok := v.(uint); ok && id != 0 {
			return id
		}
	}
	if c := ClaimsFromContext(ctx); c != nil && c.CompanyID != 0 {
		return c.CompanyID
	}
	return 0
}

func CliUserFromContext(ctx context.Context) string {
	if v := ctx.Value(ctxKeyCliUser); v != nil {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return s
		}
	}
	if c := ClaimsFromContext(ctx); c != nil && strings.TrimSpace(c.CliUser) != "" {
		return c.CliUser
	}
	return ""
}

// ----- role checks -----

// HasRequiredRole returns true if user has at least one of requiredRoles.
func HasRequiredRole(userRoles []string, requiredRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true
	}
	// normalize + fast lookup
	set := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		r = strings.TrimSpace(r)
		if r != "" {
			set[r] = struct{}{}
		}
	}
	for _, need := range requiredRoles {
		if _, ok := set[strings.TrimSpace(need)]; ok {
			return true
		}
	}
	return false
}

// splitCSV splits "ROLE_1,ROLE_2" -> []string{"ROLE_1","ROLE_2"}
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// GetMethodName extracts method name (e.g. "GetCompany") from a function value.
func GetMethodName(method interface{}) string {
	if method == nil {
		return ""
	}
	v := reflect.ValueOf(method)
	if v.Kind() != reflect.Func {
		return ""
	}
	full := runtime.FuncForPC(v.Pointer()).Name()
	// e.g. "(*company_service.CompanyService).GetCompany-fm"
	parts := strings.Split(full, ".")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, "-fm")
}

func GetAccessToken(ctx context.Context) (string, error) {
	// Try Kratos transport first (HTTP)
	if tr, ok := transport.FromServerContext(ctx); ok && tr != nil {
		h := tr.RequestHeader().Get("Authorization")
		if tok := paseto.SkipBearer(h); tok != "" {
			return tok, nil
		}
	}
	// Try gRPC metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("authorization"); len(vals) > 0 {
			if tok := paseto.SkipBearer(vals[0]); tok != "" {
				return tok, nil
			}
		}
	}
	return "", ErrMissingAuthorizationHeader
}
