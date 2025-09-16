package template_service

import (
	"context"

	template "service/api/template/v1"
)

// DeleteTemplateById delete template by id and return response with meta
func (s *TemplatesService) DeleteTemplateById(ctx context.Context, req *template.DeleteTemplateByIdRequest) (*template.DeleteTemplateByIdResponse, error) {

	// use case
	err := s.uc.DeleteTemplateById(ctx, uint(req.Id))
	if err != nil {
		return &template.DeleteTemplateByIdResponse{
			Meta: &template.MetaResponse{
				Code:    template.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	// return response
	return &template.DeleteTemplateByIdResponse{
		Meta: &template.MetaResponse{
			Code:    template.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
