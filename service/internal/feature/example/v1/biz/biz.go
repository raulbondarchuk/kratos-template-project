package example_biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ExampleRepo interface {

	FindExamples(ctx context.Context, id *uint, name *string) ([]Example, error)
	UpsertExample(ctx context.Context, in *Example) (*Example, error)
	DeleteExampleById(ctx context.Context, id uint) error
}
type ExampleUsecase struct {
	repo ExampleRepo
	log  *log.Helper
}

func NewExampleUsecase(repo ExampleRepo, logger log.Logger) *ExampleUsecase {
	return &ExampleUsecase{repo: repo, log: log.NewHelper(logger)}
}