package data

import (
	"context"
)

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type transaction struct {
	data *Data
}

func NewTransaction(data *Data) Transaction {
	return &transaction{
		data: data,
	}
}

// ExecTx executes the function in a transaction
func (t *transaction) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.data.WithTx(ctx, fn)
}

// ================================================
// ================ EJEMPLO =======================
// ================================================

// err = s.tx.ExecTx(ctx, func(ctx context.Context) error { }
// En service (NewService) tiene que existir tx -> data.Transaction en caso de que si queremos usar transacciones
