package template_repo

import (
	"context"
	"service/internal/data/model"
)

func (repo *templateRepo) DeleteTemplateById(ctx context.Context, id uint) error {
	return repo.data.DB.Delete(&model.Templates{}, id).Error
}
