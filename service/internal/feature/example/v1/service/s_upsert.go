package example_service

import (
	"context"

	api_example   "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	srvmeta      "service/internal/server/http/meta"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *ExampleService) UpsertExample(ctx context.Context, req *api_example.UpsertExampleRequest) (*api_example.UpsertExampleResponse, error) {
	in := &example_biz.Example{
		ID:   uint(req.GetId()),
		Name: req.GetName(),
	}
	res, err := s.uc.UpsertExample(ctx, in)
	if err != nil {
		return &api_example.UpsertExampleResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to upsert example"), map[string]string{"error": err.Error()}),
		}, nil
	}

	dto, err := generic.ToDTOGeneric[example_biz.Example, api_example.Example](*res)
	if err != nil {
		return &api_example.UpsertExampleResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to marshal dto"), map[string]string{"error": err.Error()}),
		}, nil
	}
	dto.CreatedAt = converter.ConvertToGoogleTimestamp(res.CreatedAt)
	dto.UpdatedAt = converter.ConvertToGoogleTimestamp(res.UpdatedAt)

	return &api_example.UpsertExampleResponse{
		Item: &dto,
		Meta: srvmeta.MetaOK("OK"),
	}, nil
}