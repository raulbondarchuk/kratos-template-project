package example_repo

import (
	"context"
	example_biz "service/internal/feature/example/v1/biz"
)

func (r *exampleRepo) UpsertExample(ctx context.Context, in *example_biz.Example) (*example_biz.Example, error) {
	return in, nil
}