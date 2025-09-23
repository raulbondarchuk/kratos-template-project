package example_biz

import "context"

func (uc *ExampleUsecase) FindExamples(ctx context.Context, id *uint, name *string) ([]Example, error) {
	return uc.repo.FindExamples(ctx, id, name)
}
func (uc *ExampleUsecase) UpsertExample(ctx context.Context, in *Example) (*Example, error) {
	return uc.repo.UpsertExample(ctx, in)
}
func (uc *ExampleUsecase) DeleteExampleById(ctx context.Context, id uint) error {
	return uc.repo.DeleteExampleById(ctx, id)
}