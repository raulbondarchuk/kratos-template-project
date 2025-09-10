package template_repo

import (
	"context"
	"service/internal/data/model"
	template_biz "service/internal/feature/template/biz"
	"service/pkg/generic"
)

// UpsertTemplate inserts new template or updates existing
func (r *templateRepo) UpsertTemplate(ctx context.Context, t *template_biz.Template) (*template_biz.Template, error) {
	var dbModel model.Templates
	var err error

	// biz -> db model
	dbModel, err = generic.ToModelGeneric[template_biz.Template, model.Templates](*t)
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
			Model(&model.Templates{}).
			Where("id = ?", dbModel.ID).
			Updates(&dbModel).Error; err != nil {
			return nil, err
		}
	}

	if err := r.data.DB.WithContext(ctx).
		Preload("Type").
		First(&dbModel, dbModel.ID).Error; err != nil {
		return nil, err
	}

	// db model -> biz
	updated, err := generic.ToDomainGeneric[model.Templates, template_biz.Template](dbModel)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
