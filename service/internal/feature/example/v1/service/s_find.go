package example_service

import (
	"context"

	api_example   "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	httperr      "service/internal/server/http/middleware/errors"
	reason       "service/internal/middleware/http_reason"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *ExampleService) FindExamples(ctx context.Context, req *api_example.FindExamplesRequest) (*api_example.FindExamplesResponse, error) {
	var idPtr *uint
	var namePtr *string

	if req.Id != nil && *req.Id != 0 {
		v := uint(*req.Id)
		idPtr = &v
	}
	if req.Name != nil && *req.Name != "" {
		v := *req.Name
		namePtr = &v
	}

	bizRes, err := s.uc.FindExamples(ctx, idPtr, namePtr)
	if err != nil {
		return nil, httperr.Internal(reason.ReasonDatabase, err.Error(), nil)
	}

	dto, err := generic.ToDTOSliceGeneric[example_biz.Example, api_example.Example](bizRes)
	if err != nil {
		return nil, httperr.Internal(reason.ReasonGeneric, err.Error(), nil)
	}
	for i := range bizRes {
		dto[i].CreatedAt = generic.ConvertToGoogleTimestamp(bizRes[i].CreatedAt)
		dto[i].UpdatedAt = generic.ConvertToGoogleTimestamp(bizRes[i].UpdatedAt)
	}

	return &api_example.FindExamplesResponse{
		Items: generic.ToPointerSliceGeneric(dto),
		Total: uint32(len(dto)),
	}, nil
}