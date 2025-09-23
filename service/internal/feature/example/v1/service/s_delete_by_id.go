package example_service

import (
	"context"

	api_example "service/api/example/v1"
	httperr    "service/internal/server/http/middleware/errors"
	reason     "service/internal/middleware/http_reason"
)

func (s *ExampleService) DeleteExampleById(ctx context.Context, req *api_example.DeleteExampleByIdRequest) (*api_example.DeleteExampleByIdResponse, error) {
	if err := s.uc.DeleteExampleById(ctx, uint(req.GetId())); err != nil {
		// if desired, you can differentiate NotFound/Conflict and etc.
		return nil, httperr.Internal(reason.ReasonDatabase, err.Error(), nil)
	}
	return &api_example.DeleteExampleByIdResponse{}, nil
}