package feature

import (

	example_v1_service "service/internal/feature/example/v1/service"
	example_v1 "service/internal/feature/example/v1"
	"service/internal/server/middleware/auth/authz/endpoint"

	"github.com/google/wire"
)

// ProvideAuthGroups get all groups of services who requires authentication and authorization
func ProvideAuthGroups(
	exampleV1Svc *example_v1_service.ExampleService,

	// Add other services there
) []endpoint.ServiceGroup {
	return []endpoint.ServiceGroup{
		example_v1.GetServiceEndpoints(exampleV1Svc),
	}
}

var ProviderAuthSet = wire.NewSet(ProvideAuthGroups)
