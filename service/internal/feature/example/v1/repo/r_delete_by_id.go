package example_repo

import (
	"context"
	"service/internal/data/model"
)

func (repo *exampleRepo) DeleteExampleById(ctx context.Context, id uint) error {
	return repo.data.DB.Delete(&model.Examples{}, id).Error
}
