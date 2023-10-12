package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type auth interface {
	Check(r *http.Request) (ctx context.Context, err error)
}

var (
	ErrNotAllowed = errors.New("not allowed")
)

var SessionValidate = func(
	auth auth,
	needAuthFunc func(w http.ResponseWriter, message string) error,
) HTTPMiddleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ctx, err := auth.Check(r)
			if err != nil {
				if !errors.Is(err, ErrNotAllowed) {
					slog.Default().Error("session validate", slog.Any("error", err))
				}

				err = needAuthFunc(w, "")
				if err != nil {
					slog.Default().Error("write need auth", slog.Any("error", err))
				}

				return
			}

			// go next if auth success
			next(w, r.WithContext(ctx), p)
		}
	}
}
