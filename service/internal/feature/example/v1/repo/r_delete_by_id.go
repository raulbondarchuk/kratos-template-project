package example_repo

import (
	"context"
	"service/internal/data/model"
)

// TODO: Mock implementation, no va a funcionar con la bdd real si tenemos mock TRUE

func (repo *exampleRepo) DeleteExampleById(ctx context.Context, id uint) error {
	if useMock {
		// Mock implementation
		if id == 0 {
			return repo.data.DB.Error
		}
		return nil
	}
	// Real implementation
	return repo.data.DB.Delete(&model.Examples{}, id).Error
}
