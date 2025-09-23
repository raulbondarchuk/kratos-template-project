package example_service

import (
	"context"

	api_example "service/api/example/v1"
	srvmeta    "service/internal/server/http/meta"
)

func (s *ExampleService) DeleteExampleById(ctx context.Context, req *api_example.DeleteExampleByIdRequest) (*api_example.DeleteExampleByIdResponse, error) {
	if err := s.uc.DeleteExampleById(ctx, uint(req.GetId())); err != nil {
		return &api_example.DeleteExampleByIdResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to delete example"), map[string]string{"error": err.Error()}),
		}, nil
	}
	return &api_example.DeleteExampleByIdResponse{
		Meta: srvmeta.MetaOK("OK"),
	}, nil
}