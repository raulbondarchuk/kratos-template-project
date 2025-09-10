package template_biz

import "context"

// ListTemplates returns all available templates
func (uc *TemplateUsecase) ListTemplates(ctx context.Context) ([]Template, error) {
	return uc.repo.ListTemplates(ctx)
}

// UpsertTemplate inserts new template or updates existing
func (uc *TemplateUsecase) UpsertTemplate(ctx context.Context, t *Template) (*Template, error) {
	updated, err := uc.repo.UpsertTemplate(ctx, t)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// DeleteTemplateById deletes a template by its ID
func (uc *TemplateUsecase) DeleteTemplateById(ctx context.Context, id uint) error {
	return uc.repo.DeleteTemplateById(ctx, id)
}

func (uc *TemplateUsecase) ReceiveTemplate(topic string, message string) error {
	uc.log.Debugf("Received message on topic %s", topic)
	return nil
}
