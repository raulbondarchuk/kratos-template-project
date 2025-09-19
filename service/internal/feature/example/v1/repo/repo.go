package example_repo

import (
	"service/internal/data"
	example_biz "service/internal/feature/example/v1/biz"

	"github.com/go-kratos/kratos/v2/log"
)

const useMock = true // Switch between mock/real implementation

// exampleRepo is a concrete implementation of biz.ExampleRepo
type exampleRepo struct {
	data *data.Data
	log  *log.Helper
}

// NewExampleRepo creates a new example repository
func NewExampleRepo(data *data.Data, logger log.Logger) example_biz.ExampleRepo {
	return &exampleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
