package migrations

import (
	"service/internal/data/model"
)

// Models to migrate
var MODELS_TO_MIGRATE = []any{
	// TODO: your models SQL
	model.Types{},
	model.Templates{},
}
