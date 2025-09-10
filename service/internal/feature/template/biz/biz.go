package template_biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// TemplateRepo describes the template repository interface
type TemplateRepo interface {
	ListTemplates(ctx context.Context) ([]Template, error)
	UpsertTemplate(ctx context.Context, t *Template) (*Template, error)
	DeleteTemplateById(ctx context.Context, id uint) error
}

// TemplateUsecase handles template business logic
type TemplateUsecase struct {
	repo TemplateRepo
	log  *log.Helper
}

// NewTemplateUsecase creates a new template usecase
func NewTemplateUsecase(repo TemplateRepo, logger log.Logger) *TemplateUsecase {
	return &TemplateUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
