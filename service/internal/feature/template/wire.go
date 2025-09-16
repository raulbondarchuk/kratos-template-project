package template

import (
	api_template "service/api/template/v1"
	template_biz "service/internal/feature/template/biz"
	template_repo "service/internal/feature/template/repo"
	template_service "service/internal/feature/template/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// your providers
	template_repo.NewTemplateRepo,
	template_biz.NewTemplateUsecase,
	template_service.NewTemplateService,

	// bind service to interfaces that buf/protoc generates
	wire.Bind(new(api_template.TemplatesServiceHTTPServer), new(*template_service.TemplatesService)),
	wire.Bind(new(api_template.TemplatesServiceServer), new(*template_service.TemplatesService)),

	// registrers (from registrars.go file)
	NewTemplatesHTTPRegistrer,
	NewTemplatesGRPCRegistrer,
)
