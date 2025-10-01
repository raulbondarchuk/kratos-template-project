package example

import (
	example_service "service/internal/feature/example/v1/service"
	"service/internal/server/middleware/auth/authz/endpoint"
)

const (
	RoleExample = "TEST1"
)

// Endpoints with required roles (versioned)
func GetExamplev1Endpoints(svc *example_service.ExampleService) endpoint.ServiceGroup {
	return endpoint.ServiceGroup{
		Name: "examplev1",
		Methods: []endpoint.ServiceMethod{
			// Examples (uncomment and replace with real service methods):
			// endpoint.NewServiceMethod(svc, svc.ListExamples),
			// endpoint.NewServiceMethod(svc, svc.UpsertExcel, RoleExample),
			// other methods...
		},
	}
}

// Backward-compatible alias (without version): calls versioned function
func GetServiceEndpoints(svc *example_service.ExampleService) endpoint.ServiceGroup {
	return GetExamplev1Endpoints(svc)
}