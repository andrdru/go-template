package ctxsess

import (
	"context"

	"github.com/andrdru/go-template/internal/entities"
)

type (
	ctxKey string
)

const (
	keyLogger ctxKey = "key"
)

// Set session to context
func Set(parent context.Context, session *entities.Session) context.Context {
	return context.WithValue(parent, keyLogger, session)
}

// Get session from context
func Get(ctx context.Context) *entities.Session {
	data := ctx.Value(keyLogger)
	if data != nil {
		if ret, ok := data.(*entities.Session); ok {
			return ret
		}
	}

	return nil
}
