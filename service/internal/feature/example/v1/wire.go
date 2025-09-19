package example

import (
	api_example "service/api/example/v1"
	biz "service/internal/feature/example/v1/biz"
	repo "service/internal/feature/example/v1/repo"
	service "service/internal/feature/example/v1/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// your providers
	repo.NewExampleRepo,
	biz.NewExampleUsecase,
	service.NewExampleService,

	// bind service to interfaces that buf/protoc generates
	wire.Bind(new(api_example.Examplev1ServiceHTTPServer), new(*service.ExampleService)),
	wire.Bind(new(api_example.Examplev1ServiceServer), new(*service.ExampleService)),

	// registrers (from registrers.go file)
	NewExampleHTTPRegistrer,
	NewExampleGRPCRegistrer,
)
