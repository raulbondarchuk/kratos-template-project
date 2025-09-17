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

	wire.Bind(new(api_prueba.PruebaServiceHTTPServer), new(*prueba_service.PruebaService)),
	wire.Bind(new(api_prueba.PruebaServiceServer),     new(*prueba_service.PruebaService)),

	NewPruebaHTTPRegister,
	NewPruebaGRPCRegister,
)