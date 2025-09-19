package example_repo

import (
	"context"
	"service/internal/data/model"
	example_biz "service/internal/feature/example/v1/biz"
	"service/pkg/generic"
)

// ListTemplates returns all templates
func (r *exampleRepo) ListExamples(ctx context.Context) ([]example_biz.Example, error) {

	var examples []model.Examples
	if err := r.data.DB.Preload("TypeExamples").Find(&examples).Error; err != nil {
		return nil, err
	}

	return generic.ToDomainSliceGeneric[model.Examples, example_biz.Example](examples)
}
