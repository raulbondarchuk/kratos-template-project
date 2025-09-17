package prueba_service

import (
	"context"

	api_prueba "service/api/prueba/v1"
)

func (s *PruebaService) DeletePruebaById(ctx context.Context, req *api_prueba.DeletePruebaByIdRequest) (*api_prueba.DeletePruebaByIdResponse, error) {
	if err := s.uc.DeletePruebaById(ctx, uint(req.Id)); err != nil {
		return &api_prueba.DeletePruebaByIdResponse{
			Meta: &api_prueba.MetaResponse{
				Code:    api_prueba.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}
	return &api_prueba.DeletePruebaByIdResponse{
		Meta: &api_prueba.MetaResponse{
			Code:    api_prueba.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}