package template_repo

import (
	"service/internal/data"
	template_biz "service/internal/feature/template/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// templateRepo is a concrete implementation of biz.TemplateRepo
type templateRepo struct {
	data *data.Data
	log  *log.Helper
}

// NewTemplateRepo creates a new template repository
func NewTemplateRepo(data *data.Data, logger log.Logger) template_biz.TemplateRepo {
	return &templateRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
