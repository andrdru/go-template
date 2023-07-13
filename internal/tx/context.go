package tx

import (
	"context"
	"database/sql"
)

type (
	ctxKey string
)

const (
	keyTransaction ctxKey = "key"
)

// ctxSetTx tx to context
func ctxSetTx(parent context.Context, item *sql.Tx) context.Context {
	return context.WithValue(parent, keyTransaction, item)
}

// ctxGetTx tx from context
func ctxGetTx(ctx context.Context) *sql.Tx {
	data := ctx.Value(keyTransaction)
	if data != nil {
		if ret, ok := data.(*sql.Tx); ok {
			return ret
		}
	}

	return nil
}
