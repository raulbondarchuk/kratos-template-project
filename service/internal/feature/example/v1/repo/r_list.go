package example_repo

import (
	"context"
	example_biz "service/internal/feature/example/v1/biz"
)

func (r *exampleRepo) FindExamples(ctx context.Context, id *uint, name *string) ([]example_biz.Example, error) {
	return []example_biz.Example{}, nil
}