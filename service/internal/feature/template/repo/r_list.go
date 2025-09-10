package template_repo

import (
	"context"
	"service/internal/data/model"
	template_biz "service/internal/feature/template/biz"
	"service/pkg/generic"
)

// ListTemplates returns all templates
func (r *templateRepo) ListTemplates(ctx context.Context) ([]template_biz.Template, error) {

	var templates []model.Templates
	if err := r.data.DB.Preload("Type").Find(&templates).Error; err != nil {
		return nil, err
	}

	return generic.ToDomainSliceGeneric[model.Templates, template_biz.Template](templates)
}
