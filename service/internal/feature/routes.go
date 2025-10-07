package feature

import (

	"service/internal/server/middleware/auth/authz/endpoint"

	"github.com/google/wire"
)

// ProvideAuthGroups get all groups of services who requires authentication and authorization
func ProvideAuthGroups(

	// Add other services there
) []endpoint.ServiceGroup {
	return []endpoint.ServiceGroup{}
}

var ProviderAuthSet = wire.NewSet(ProvideAuthGroups)
