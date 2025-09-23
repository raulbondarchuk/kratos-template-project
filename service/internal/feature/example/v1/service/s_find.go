package example_service

import (
	"context"

	api_example   "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	srvmeta      "service/internal/server/http/meta"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *ExampleService) FindExamples(ctx context.Context, req *api_example.FindExamplesRequest) (*api_example.FindExamplesResponse, error) {
	// presence-aware (optional fields)
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
		return &api_example.FindExamplesResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to find example"), map[string]string{"error": err.Error()}),
		}, nil
	}

	if len(bizRes) == 0 {
		return &api_example.FindExamplesResponse{
			Meta: srvmeta.MetaNoContent("no items"),
		}, nil
	}

	dto, err := generic.ToDTOSliceGeneric[example_biz.Example, api_example.Example](bizRes)
	if err != nil {
		return &api_example.FindExamplesResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to marshal dto"), map[string]string{"error": err.Error()}),
		}, nil
	}

	for i := range bizRes {
		dto[i].CreatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].CreatedAt)
		dto[i].UpdatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].UpdatedAt)
	}

	return &api_example.FindExamplesResponse{
		Items: generic.ToPointerSliceGeneric(dto),
		Meta:  srvmeta.MetaOK("OK"),
	}, nil
}