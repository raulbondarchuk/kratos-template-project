package template_service

import (
	"context"

	template "service/api/template/v1"
	template_biz "service/internal/feature/template/biz"
	"service/pkg/converter"
	"service/pkg/generic"
)

// UpsertTemplate upsert template and return response with template.
// If id is 0, it will create a new template, otherwise it will update the existing template.
func (s *TemplatesService) UpsertTemplate(ctx context.Context, req *template.UpsertTemplateRequest) (*template.UpsertTemplateResponse, error) {

	// biz-level model
	request := &template_biz.Template{
		ID:   uint(req.Id),
		Name: req.Name,
	}

	// use case
	bizResult, err := s.uc.UpsertTemplate(ctx, request)
	if err != nil {
		return &template.UpsertTemplateResponse{
			Meta: &template.MetaResponse{
				Code:    template.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// ✨ use generic for mapping Biz -> DTO
	protoResult, err := generic.ToDTOGeneric[template_biz.Template, template.Template](*bizResult)
	if err != nil {
		return &template.UpsertTemplateResponse{
			Meta: &template.MetaResponse{
				Code:    template.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: "Mapping error",
			},
		}, nil
	}

	// fill timestamps
	protoResult.CreatedAt = converter.ConvertToGoogleTimestamp(bizResult.CreatedAt)
	protoResult.UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult.UpdatedAt)
	protoResult.Type.CreatedAt = converter.ConvertToGoogleTimestamp(bizResult.Type.CreatedAt)
	protoResult.Type.UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult.Type.UpdatedAt)

	return &template.UpsertTemplateResponse{
		Template: &protoResult,
		Meta: &template.MetaResponse{
			Code:    template.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
