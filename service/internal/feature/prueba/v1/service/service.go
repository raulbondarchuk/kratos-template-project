package prueba_service

import (
	api_prueba "service/api/prueba/v1"
	prueba_biz "service/internal/feature/prueba/v1/biz"
)

type PruebaService struct {
	api_prueba.UnimplementedPruebav1ServiceServer
	uc *prueba_biz.PruebaUsecase
}

func NewPruebaService(uc *prueba_biz.PruebaUsecase) *PruebaService {
	return &PruebaService{uc: uc}
}