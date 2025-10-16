package example_biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type ExampleRepo interface {
}
type ExampleUsecase struct {
	repo ExampleRepo
	log  *log.Helper
}

func NewExampleUsecase(repo ExampleRepo, logger log.Logger) *ExampleUsecase {
	return &ExampleUsecase{repo: repo, log: log.NewHelper(logger)}
}