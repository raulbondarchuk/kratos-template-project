package prueba_service

import (
	"context"

	common      "service/api/common/v1"
	api_prueba  "service/api/prueba/v1"
	prueba_biz "service/internal/feature/prueba/v1/biz"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *PruebaService) UpsertPrueba(ctx context.Context, req *api_prueba.UpsertPruebaRequest) (*api_prueba.UpsertPruebaResponse, error) {
	in := &prueba_biz.Prueba{
		ID:   uint(req.GetId()),
		Name: req.GetName(),
	}
	res, err := s.uc.UpsertPrueba(ctx, in)
	if err != nil {
		return &api_prueba.UpsertPruebaResponse{
			Meta: &common.MetaResponse{
				Code:    common.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: "failed to upsert prueba",
				Details: map[string]string{"error": err.Error()},
			},
		}, nil
	}

	dto, err := generic.ToDTOGeneric[prueba_biz.Prueba, api_prueba.Prueba](*res)
	if err != nil {
		return &api_prueba.UpsertPruebaResponse{
			Meta: &common.MetaResponse{
				Code:    common.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: "failed to marshal dto",
				Details: map[string]string{"error": err.Error()},
			},
		}, nil
	}
	dto.CreatedAt = converter.ConvertToGoogleTimestamp(res.CreatedAt)
	dto.UpdatedAt = converter.ConvertToGoogleTimestamp(res.UpdatedAt)

	return &api_prueba.UpsertPruebaResponse{
		Item: &dto,
		Meta: &common.MetaResponse{
			Code:    common.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}