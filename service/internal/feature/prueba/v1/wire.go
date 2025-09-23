package prueba

import (
	api_prueba "service/api/prueba/v1"
	prueba_biz "service/internal/feature/prueba/v1/biz"
	prueba_repo "service/internal/feature/prueba/v1/repo"
	prueba_service "service/internal/feature/prueba/v1/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	prueba_repo.NewPruebaRepo,
	prueba_biz.NewPruebaUsecase,
	prueba_service.NewPruebaService,

	// map generated service interfaces (versioned) to our implementation
	wire.Bind(new(api_prueba.Pruebav1ServiceHTTPServer), new(*prueba_service.PruebaService)),
	wire.Bind(new(api_prueba.Pruebav1ServiceServer),     new(*prueba_service.PruebaService)),

	// module-local registrars
	NewPruebaHTTPRegistrer,
	NewPruebaGRPCRegistrer,
)