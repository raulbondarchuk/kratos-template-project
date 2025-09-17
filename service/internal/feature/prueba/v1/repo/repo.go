package prueba_repo

import (
	"service/internal/data"
	prueba_biz "service/internal/feature/prueba/v1/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type pruebaRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewPruebaRepo(data *data.Data, logger log.Logger) prueba_biz.PruebaRepo {
	return &pruebaRepo{data: data, log: log.NewHelper(logger)}
}