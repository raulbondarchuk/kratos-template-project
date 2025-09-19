package example_biz

import "context"

// ListTemplates returns all available templates
func (uc *ExampleUsecase) ListExamples(ctx context.Context) ([]Example, error) {
	return uc.repo.ListExamples(ctx)
}

// UpsertTemplate inserts new template or updates existing
func (uc *ExampleUsecase) UpsertExample(ctx context.Context, t *Example) (*Example, error) {
	updated, err := uc.repo.UpsertExample(ctx, t)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// DeleteExampleById deletes a template by its ID
func (uc *ExampleUsecase) DeleteExampleById(ctx context.Context, id uint) error {
	return uc.repo.DeleteExampleById(ctx, id)
}

func (uc *ExampleUsecase) ReceiveExample(topic string, message string) error {
	uc.log.Debugf("Received message on topic %s", topic)
	return nil
}
