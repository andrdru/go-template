package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type auth interface {
	Check(r *http.Request) (ctx context.Context, err error)
}

var (
	ErrNotAllowed = errors.New("not allowed")
)

var SessionValidate = func(
	logger zerolog.Logger,
	auth auth,
	needAuthFunc func(w http.ResponseWriter, message string) error,
) HTTPMiddleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ctx, err := auth.Check(r)
			if err != nil {
				if !errors.Is(err, ErrNotAllowed) {
					logger.Error().Err(err).Msgf("session validate")
				}

				err = needAuthFunc(w, "")
				if err != nil {
					logger.Error().Err(err).Msgf("write need auth")
				}

				return
			}

			// go next if auth success
			next(w, r.WithContext(ctx), p)
		}
	}
}
