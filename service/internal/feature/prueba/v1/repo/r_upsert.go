package prueba_repo

import (
	"context"
	prueba_biz "service/internal/feature/prueba/v1/biz"
)

func (r *pruebaRepo) UpsertPrueba(ctx context.Context, in *prueba_biz.Prueba) (*prueba_biz.Prueba, error) {
	return in, nil
}