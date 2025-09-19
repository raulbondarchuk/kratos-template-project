package example_service

import (
	"context"

	example "service/api/example/v1"
)

// DeleteTemplateById delete template by id and return response with meta
func (s *ExampleService) DeleteExampleById(ctx context.Context, req *example.DeleteExampleByIdRequest) (*example.DeleteExampleByIdResponse, error) {

	// use case
	err := s.uc.DeleteExampleById(ctx, uint(req.Id))
	if err != nil {
		return &example.DeleteExampleByIdResponse{
			Meta: &example.MetaResponse{
				Code:    example.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// return response
	return &example.DeleteExampleByIdResponse{
		Meta: &example.MetaResponse{
			Code:    example.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
