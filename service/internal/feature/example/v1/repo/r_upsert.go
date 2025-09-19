package example_repo

import (
	"context"
	"service/internal/data/model"
	example_biz "service/internal/feature/example/v1/biz"
	"service/pkg/generic"
)

// UpsertTemplate inserts new template or updates existing
func (r *exampleRepo) UpsertExample(ctx context.Context, t *example_biz.Example) (*example_biz.Example, error) {
	var dbModel model.Examples
	var err error

	// biz -> db model
	dbModel, err = generic.ToModelGeneric[example_biz.Example, model.Examples](*t)
	if err != nil {
		return nil, err
	}

	if dbModel.ID == 0 {
		if err := r.data.DB.WithContext(ctx).Create(&dbModel).Error; err != nil {
			return nil, err
		}
	} else {
		// update
		if err := r.data.DB.WithContext(ctx).
			Model(&model.Examples{}).
			Where("id = ?", dbModel.ID).
			Updates(&dbModel).Error; err != nil {
			return nil, err
		}
	}

	if err := r.data.DB.WithContext(ctx).
		Preload("TypeExamples").
		First(&dbModel, dbModel.ID).Error; err != nil {
		return nil, err
	}

	// db model -> biz
	updated, err := generic.ToDomainGeneric[model.Examples, example_biz.Example](dbModel)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
