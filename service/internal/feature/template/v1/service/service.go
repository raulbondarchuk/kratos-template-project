package template_service

import (
	api_template "service/api/template/v1"
	template_biz "service/internal/feature/template/v1/biz"
)

// TemplateService implements the template service
type TemplateService struct {
	api_template.UnimplementedTemplatev1ServiceServer

	uc *template_biz.TemplateUsecase
}

// NewTemplateService creates a new template service
func NewTemplateService(uc *template_biz.TemplateUsecase) *TemplateService {
	return &TemplateService{uc: uc}
}
