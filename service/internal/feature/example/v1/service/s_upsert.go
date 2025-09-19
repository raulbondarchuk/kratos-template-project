package example_service

import (
	"context"

	example "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	"service/pkg/converter"
	"service/pkg/generic"
)

// UpsertTemplate upsert template and return response with template.
// If id is 0, it will create a new template, otherwise it will update the existing template.
func (s *ExampleService) UpsertExample(ctx context.Context, req *example.UpsertExampleRequest) (*example.UpsertExampleResponse, error) {

	// biz-level model
	request := &example_biz.Example{
		ID:   uint(req.Id),
		Name: req.Name,
	}

	// use case
	bizResult, err := s.uc.UpsertExample(ctx, request)
	if err != nil {
		return &example.UpsertExampleResponse{
			Meta: &example.MetaResponse{
				Code:    example.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// ✨ use generic for mapping Biz -> DTO
	protoResult, err := generic.ToDTOGeneric[example_biz.Example, example.Example](*bizResult)
	if err != nil {
		return &example.UpsertExampleResponse{
			Meta: &example.MetaResponse{
				Code:    example.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: "Mapping error",
			},
		}, nil
	}

	// fill timestamps
	protoResult.CreatedAt = converter.ConvertToGoogleTimestamp(bizResult.CreatedAt)
	protoResult.UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult.UpdatedAt)
	protoResult.Type.CreatedAt = converter.ConvertToGoogleTimestamp(bizResult.Type.CreatedAt)
	protoResult.Type.UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult.Type.UpdatedAt)

	return &example.UpsertExampleResponse{
		Example: &protoResult,
		Meta: &example.MetaResponse{
			Code:    example.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
