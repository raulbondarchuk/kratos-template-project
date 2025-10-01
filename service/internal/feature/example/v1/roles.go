package example

import (
	example_service "service/internal/feature/example/v1/service"
	"service/internal/server/middleware/auth/authz/endpoint"
)

const (
	RoleTEST1 = "TEST1"
)

// Endpoints with required roles
func GetServiceEndpoints(svc *example_service.ExampleService) endpoint.ServiceGroup {
	return endpoint.ServiceGroup{
		Name: "example",
		Methods: []endpoint.ServiceMethod{
			// Examples (uncomment and replace with real service methods):
			// endpoint.NewServiceMethod(svc, svc.ListExamples),
			// endpoint.NewServiceMethod(svc, svc.UpsertExcel, RoleTEST1),
			// other methods...
		},
	}
}