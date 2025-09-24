package data

import (
	"context"

	"gorm.io/gorm"
)

// contextTxKey is a special type for the transaction key in context
type contextTxKey struct{}

// txKey is the single instance of the transaction key
var txKey = contextTxKey{}

// DB returns the transaction from the context if it exists, otherwise returns the normal connection to the DB
func (d *Data) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return d.db
}

// WithTx wraps the function in a transaction if it doesn't exist in the context
func (d *Data) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// If transaction already exists, use it
	if _, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return fn(ctx)
	}

	// If there is no transaction, create a new one
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, txKey, tx)
		return fn(ctx)
	})
}
