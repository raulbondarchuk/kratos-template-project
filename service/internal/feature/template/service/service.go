package template_service

import (
	template "service/api/template"
	template_biz "service/internal/feature/template/biz"
)

// TemplateService implements the template service
type TemplatesService struct {
	template.UnimplementedTemplatesServer

	uc *template_biz.TemplateUsecase
}

// NewTemplateService creates a new template service
func NewTemplateService(uc *template_biz.TemplateUsecase) *TemplatesService {
	return &TemplatesService{uc: uc}
}
