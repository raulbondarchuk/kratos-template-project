package template

import (
	template_biz "service/internal/feature/template/biz"
	template_repo "service/internal/feature/template/repo"
	template_service "service/internal/feature/template/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	template_repo.NewTemplateRepo,
	template_biz.NewTemplateUsecase,
	template_service.NewTemplateService,
)
