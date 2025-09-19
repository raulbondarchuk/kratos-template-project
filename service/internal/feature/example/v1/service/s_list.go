package example_service

import (
	"context"

	example "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	"service/pkg/converter"
	"service/pkg/generic"
)

// ListTemplates list templates and return response with meta
func (s *ExampleService) ListExamples(ctx context.Context, req *example.ListExamplesRequest) (*example.ListExamplesResponse, error) {

	bizResult, err := s.uc.ListExamples(ctx)
	if err != nil {
		return &example.ListExamplesResponse{
			Meta: &example.MetaResponse{
				Code:    example.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	if len(bizResult) == 0 {
		return &example.ListExamplesResponse{
			Meta: &example.MetaResponse{
				Code:    example.ResponseCode_RESPONSE_CODE_OK,
				Message: "No templates found",
			},
		}, nil
	}

	protoResult, err := generic.ToDTOSliceGeneric[example_biz.Example, example.Example](bizResult)
	if err != nil {
		return &example.ListExamplesResponse{
			Meta: &example.MetaResponse{
				Code:    example.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// fill timestamps
	for index := range bizResult {
		protoResult[index].CreatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].CreatedAt)
		protoResult[index].UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].UpdatedAt)
		protoResult[index].Type.CreatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].Type.CreatedAt)
		protoResult[index].Type.UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].Type.UpdatedAt)
	}

	return &example.ListExamplesResponse{
		Examples: generic.ToPointerSliceGeneric(protoResult),
		Meta: &example.MetaResponse{
			Code:    example.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
