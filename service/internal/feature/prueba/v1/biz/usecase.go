package prueba_biz

import "context"

func (uc *PruebaUsecase) FindPruebas(ctx context.Context, id *uint, name *string) ([]Prueba, error) {
	return uc.repo.FindPruebas(ctx, id, name)
}
func (uc *PruebaUsecase) UpsertPrueba(ctx context.Context, in *Prueba) (*Prueba, error) {
	return uc.repo.UpsertPrueba(ctx, in)
}