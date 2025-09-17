package template

import (
	api_template "service/api/template/v1"
	template_biz "service/internal/feature/template/v1/biz"
	template_repo "service/internal/feature/template/v1/repo"
	template_service "service/internal/feature/template/v1/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// your providers
	template_repo.NewTemplateRepo,
	template_biz.NewTemplateUsecase,
	template_service.NewTemplateService,

	// bind service to interfaces that buf/protoc generates
	wire.Bind(new(api_template.Templatev1ServiceHTTPServer), new(*template_service.TemplateService)),
	wire.Bind(new(api_template.Templatev1ServiceServer), new(*template_service.TemplateService)),

	// registrers (from registrars.go file)
	NewTemplatesHTTPRegistrer,
	NewTemplatesGRPCRegistrer,
)
