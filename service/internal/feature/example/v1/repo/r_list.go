package example_repo

import (
	"context"
	"service/internal/data/model"
	example_biz "service/internal/feature/example/v1/biz"
	"service/pkg/generic"
	"time"
)

// TODO: Mock implementation, no va a funcionar con la bdd real si tenemos mock TRUE

// ListExamples returns all examples
func (r *exampleRepo) ListExamples(ctx context.Context) ([]example_biz.Example, error) {
	if useMock {
		// Mock data
		mockExamples := []model.Examples{
			{
				Base: model.Base{ID: 1},
				TypeExamples: model.TypesExamples{
					Base: model.Base{ID: 1},
					Name: "MOCK_TYPE_1",
				},
				TypeExamplesID: 1,
				Name:           "MOCK_EXAMPLE_1",
				UpdatedAt:      time.Now(),
			},
			{
				Base: model.Base{ID: 2},
				TypeExamples: model.TypesExamples{
					Base: model.Base{ID: 2},
					Name: "MOCK_TYPE_2",
				},
				TypeExamplesID: 2,
				Name:           "MOCK_EXAMPLE_2",
				UpdatedAt:      time.Now(),
			},
		}
		return generic.ToDomainSliceGeneric[model.Examples, example_biz.Example](mockExamples)
	}

	// Real implementation
	var examples []model.Examples
	if err := r.data.DB.Preload("TypeExamples").Find(&examples).Error; err != nil {
		return nil, err
	}
	return generic.ToDomainSliceGeneric[model.Examples, example_biz.Example](examples)
}
