package example_biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// ExampleRepo describes the example repository interface
type ExampleRepo interface {
	ListExamples(ctx context.Context) ([]Example, error)
	UpsertExample(ctx context.Context, t *Example) (*Example, error)
	DeleteExampleById(ctx context.Context, id uint) error
}

// ExampleUsecase handles example business logic
type ExampleUsecase struct {
	repo ExampleRepo
	log  *log.Helper
}

// NewExampleUsecase creates a new example usecase
func NewExampleUsecase(repo ExampleRepo, logger log.Logger) *ExampleUsecase {
	return &ExampleUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
