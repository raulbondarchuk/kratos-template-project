package template_service

import (
	"context"
	errors "service/api/errors/v1"
	template "service/bin/proto/endpoints/template"
	template_biz "service/internal/feature/template/biz"
	"service/pkg/converter"
	"service/pkg/generic"
)

// ListTemplates list templates and return response with meta
func (s *TemplatesService) ListTemplates(ctx context.Context, req *template.ListTemplatesRequest) (*template.ListTemplatesResponse, error) {

	bizResult, err := s.uc.ListTemplates(ctx)
	if err != nil {
		return &template.ListTemplatesResponse{
			Meta: &errors.StandardResponse{
				Code:    errors.ResponseCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	if len(bizResult) == 0 {
		return &template.ListTemplatesResponse{
			Meta: &errors.StandardResponse{
				Code:    errors.ResponseCode_OK,
				Message: "No templates found",
			},
		}, nil
	}

	protoResult, err := generic.ToDTOSliceGeneric[template_biz.Template, template.Template](bizResult)
	if err != nil {
		return &template.ListTemplatesResponse{
			Meta: &errors.StandardResponse{
				Code:    errors.ResponseCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// fill timestamps
	for index := range bizResult {
		protoResult[index].CreatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].CreatedAt)
		protoResult[index].UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].UpdatedAt)
		protoResult[index].Type.CreatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].Type.CreatedAt)
		protoResult[index].Type.UpdatedAt = converter.ConvertToGoogleTimestamp(bizResult[index].Type.UpdatedAt)
	}

	return &template.ListTemplatesResponse{
		Templates: generic.ToPointerSliceGeneric(protoResult),
		Meta: &errors.StandardResponse{
			Code:    errors.ResponseCode_OK,
			Message: "OK",
		},
	}, nil
}
