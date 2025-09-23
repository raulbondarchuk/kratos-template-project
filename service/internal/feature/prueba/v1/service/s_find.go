package prueba_service

import (
	"context"

	api_prueba   "service/api/prueba/v1"
	prueba_biz "service/internal/feature/prueba/v1/biz"
	srvmeta      "service/internal/server/http/meta"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *PruebaService) FindPruebas(ctx context.Context, req *api_prueba.FindPruebasRequest) (*api_prueba.FindPruebasResponse, error) {
	// presence-aware (optional fields)
	var idPtr *uint
	var namePtr *string

	if req.Id != nil && *req.Id != 0 {
	  v := uint(*req.Id)
	  idPtr = &v
	}
	if req.Name != nil && *req.Name != "" {
	  v := *req.Name
	  namePtr = &v
	}

	bizRes, err := s.uc.FindPruebas(ctx, idPtr, namePtr)
	if err != nil {
		return &api_prueba.FindPruebasResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to find prueba"), map[string]string{"error": err.Error()}),
		}, nil
	}

	if len(bizRes) == 0 {
		return &api_prueba.FindPruebasResponse{
			Meta: srvmeta.MetaNoContent("no items"),
		}, nil
	}

	dto, err := generic.ToDTOSliceGeneric[prueba_biz.Prueba, api_prueba.Prueba](bizRes)
	if err != nil {
		return &api_prueba.FindPruebasResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to marshal dto"), map[string]string{"error": err.Error()}),
		}, nil
	}

	for i := range bizRes {
		dto[i].CreatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].CreatedAt)
		dto[i].UpdatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].UpdatedAt)
	}

	return &api_prueba.FindPruebasResponse{
		Items: generic.ToPointerSliceGeneric(dto),
		Meta:  srvmeta.MetaOK("OK"),
	}, nil
}