package tx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type (
	// QueryExecutor query executor, sql.DB or sql.TX
	QueryExecutor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	// TX sql helper
	// detect query executor into repo with DB()
	// process transaction with TX()
	TX struct {
		db *sql.DB
	}

	options struct {
		level sql.IsolationLevel
	}

	Option func(*options)
)

var (
	// ErrTxOpenAlready .
	ErrTxOpenAlready = errors.New("tx open already")
)

// NewTX .
func NewTX(db *sql.DB) *TX {
	return &TX{
		db: db,
	}
}

// DB ctxGetTx query executor by context
func (T *TX) DB(ctx context.Context) QueryExecutor {
	db := ctxGetTx(ctx)
	if db != nil {
		return db
	}

	return T.db
}

// TX abstract logic from transaction details
// important: processor have to use DB() calls for properly transaction handling
func (T *TX) TX(ctx context.Context, processor func(txCtx context.Context) error, opts ...Option) error {
	if ctxGetTx(ctx) != nil {
		return ErrTxOpenAlready
	}

	var args = &options{
		level: sql.LevelDefault,
	}

	for _, opt := range opts {
		opt(args)
	}

	tx, err := T.db.BeginTx(ctx, &sql.TxOptions{Isolation: args.level})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txCtx := ctxSetTx(ctx, tx)

	err = processor(txCtx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = fmt.Errorf("rollback failed %s: %w", errRollback.Error(), err)
		}

		return fmt.Errorf("process tx: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = fmt.Errorf("rollback failed %s: %w", errRollback.Error(), err)
		}

		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func WithIsolation(level sql.IsolationLevel) Option {
	return func(args *options) {
		args.level = level
	}
}
