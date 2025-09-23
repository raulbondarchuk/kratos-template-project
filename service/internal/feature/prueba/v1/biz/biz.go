package prueba_biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type PruebaRepo interface {

	FindPruebas(ctx context.Context, id *uint, name *string) ([]Prueba, error)
	UpsertPrueba(ctx context.Context, in *Prueba) (*Prueba, error)
}
type PruebaUsecase struct {
	repo PruebaRepo
	log  *log.Helper
}

func NewPruebaUsecase(repo PruebaRepo, logger log.Logger) *PruebaUsecase {
	return &PruebaUsecase{repo: repo, log: log.NewHelper(logger)}
}