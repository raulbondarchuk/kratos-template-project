package template_service

import (
	"context"
	"service/api/errors"
	template "service/api/template"
)

// DeleteTemplateById delete template by id and return response with meta
func (s *TemplatesService) DeleteTemplateById(ctx context.Context, req *template.DeleteTemplateRequest) (*template.DeleteTemplateResponse, error) {

	// use case
	err := s.uc.DeleteTemplateById(ctx, uint(req.Id))
	if err != nil {
		return &template.DeleteTemplateResponse{
			Meta: &errors.StandardResponse{
				Code:    errors.ResponseCode_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// return response
	return &template.DeleteTemplateResponse{
		Meta: &errors.StandardResponse{
			Code:    errors.ResponseCode_OK,
			Message: "OK",
		},
	}, nil
}
