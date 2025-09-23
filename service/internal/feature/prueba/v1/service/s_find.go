package prueba_service

import (
	"context"

	common      "service/api/common/v1"
	api_prueba  "service/api/prueba/v1"
	prueba_biz "service/internal/feature/prueba/v1/biz"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *PruebaService) FindPruebas(ctx context.Context, req *api_prueba.FindPruebasRequest) (*api_prueba.FindPruebasResponse, error) {
	// presence-aware (optional fields)
	var idPtr *uint
	var namePtr *string

	if req.Id != nil && *req.Id != 0 {
		tmp := uint(*req.Id)
		idPtr = &tmp
	}
	if req.Name != nil && *req.Name != "" {
		tmp := *req.Name
		namePtr = &tmp
	}

	bizRes, err := s.uc.FindPruebas(ctx, idPtr, namePtr)
	if err != nil {
		return &api_prueba.FindPruebasResponse{
			Meta: &common.MetaResponse{
				Code:    common.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: "failed to find prueba",
				Details: map[string]string{"error": err.Error()},
			},
		}, nil
	}

	if len(bizRes) == 0 {
		return &api_prueba.FindPruebasResponse{
			Meta: &common.MetaResponse{
				Code:    common.ResponseCode_RESPONSE_CODE_NO_CONTENT,
				Message: "no items",
			},
		}, nil
	}

	dto, err := generic.ToDTOSliceGeneric[prueba_biz.Prueba, api_prueba.Prueba](bizRes)
	if err != nil {
		return &api_prueba.FindPruebasResponse{
			Meta: &common.MetaResponse{
				Code:    common.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: "failed to marshal dto",
				Details: map[string]string{"error": err.Error()},
			},
		}, nil
	}

	for i := range bizRes {
		dto[i].CreatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].CreatedAt)
		dto[i].UpdatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].UpdatedAt)
	}

	return &api_prueba.FindPruebasResponse{
		Items: generic.ToPointerSliceGeneric(dto),
		Meta: &common.MetaResponse{
			Code:    common.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}