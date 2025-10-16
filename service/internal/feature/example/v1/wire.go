package example

import (
	api_example "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	example_repo "service/internal/feature/example/v1/repo"
	example_service "service/internal/feature/example/v1/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	example_repo.NewExampleRepo,
	example_biz.NewExampleUsecase,
	example_service.NewExampleService,

	// map generated service interfaces (versioned) to our implementation
	wire.Bind(new(api_example.Examplev1ServiceHTTPServer), new(*example_service.ExampleService)),
	wire.Bind(new(api_example.Examplev1ServiceServer),     new(*example_service.ExampleService)),

	// module-local registrars
	NewExampleHTTPRegistrer,
	NewExampleGRPCRegistrer,
)