package endpoint

import (
	"context"
	"reflect"
	"runtime"
	"strings"
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
	ctxKeyRoles  ctxKey = "roles"
	ctxKeyClaims ctxKey = "claims"
)

// RolesFromContext returns roles previously stored by middleware.
func RolesFromContext(ctx context.Context) []string {
	if v := ctx.Value(ctxKeyRoles); v != nil {
		if roles, ok := v.([]string); ok {
			return roles
		}
	}
	return nil
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
