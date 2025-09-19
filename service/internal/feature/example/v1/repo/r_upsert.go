package example_repo

import (
	"context"
	"service/internal/data/model"
	example_biz "service/internal/feature/example/v1/biz"
	"service/pkg/generic"
	"time"
)

// TODO: Mock implementation, no va a funcionar con la bdd real si tenemos mock TRUE

// UpsertExample inserts new template or updates existing
func (r *exampleRepo) UpsertExample(ctx context.Context, t *example_biz.Example) (*example_biz.Example, error) {
	if useMock {
		// Mock implementation
		dbModel := model.Examples{
			Base: model.Base{
				ID: t.ID,
			},
			TypeExamples: model.TypesExamples{
				Base: model.Base{ID: t.Type.ID},
				Name: "MOCK_TYPE_AUTO",
			},
			TypeExamplesID: t.Type.ID,
			Name:           t.Name,
			UpdatedAt:      time.Now(),
		}

		// Simulate ID generation for new records
		if dbModel.ID == 0 {
			dbModel.ID = 999 // Mock new ID
		}

		// Convert back to domain model
		updated, err := generic.ToDomainGeneric[model.Examples, example_biz.Example](dbModel)
		if err != nil {
			return nil, err
		}

		return &updated, nil
	}

	// Real implementation
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
