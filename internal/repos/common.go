package repos

import (
	"context"

	"github.com/andrdru/go-template/tx"
)

type (
	// Transactor sql helper
	transactor interface {
		DB(ctx context.Context) (db tx.QueryExecutor)
		TX(ctx context.Context, processor func(txCtx context.Context) error, opts ...tx.Option) error
	}
)
