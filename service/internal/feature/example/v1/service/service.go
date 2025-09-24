package example_service

import (
	api_example "service/api/example/v1"
	example_biz "service/internal/feature/example/v1/biz"
	appdata     "service/internal/data"
)

type ExampleService struct {
	api_example.UnimplementedExamplev1ServiceServer
	uc *example_biz.ExampleUsecase
	tx appdata.Transaction
}

func NewExampleService(uc *example_biz.ExampleUsecase, tx appdata.Transaction) *ExampleService {
	return &ExampleService{uc: uc, tx: tx}
}