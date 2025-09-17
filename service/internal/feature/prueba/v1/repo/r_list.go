package prueba_repo

import (
	"context"
	prueba_biz "service/internal/feature/prueba/v1/biz"
)

func (r *pruebaRepo) FindPruebas(ctx context.Context, id *uint, name *string) ([]prueba_biz.Prueba, error) {
	return []prueba_biz.Prueba{}, nil
}