package app

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type (
	//closer graceful shutdown handler
	closer func(ctx context.Context) (description string, err error)

	bootstrap struct {
		closers []closer
	}
)

func Run(logger zerolog.Logger) (code int) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info().Msg("app staring")
	defer logger.Info().Msg("app finished")

	boot, err := initApp(logger)
	if err != nil {
		logger.Error().Err(err).Msg("init app")
		return 1
	}

	defer func() {
		ctxCloser, cancelCloser := context.WithTimeout(context.Background(), 15*time.Second)
		defer func() {
			cancelCloser()
		}()

		boot.shutdown(ctxCloser, logger)
	}()

	logger.Info().Msg("app started successfully")
	<-ctx.Done()

	return 0
}

func initApp(logger zerolog.Logger) (boot bootstrap, err error) {
	// todo
	return bootstrap{}, nil
}

// shutdown close app resources
func (b *bootstrap) shutdown(ctx context.Context, logger zerolog.Logger) {
	for _, f := range b.closers {
		description, errClose := f(ctx)
		if errClose != nil {
			logger.Err(errClose).Msgf("close: %s", description)
			continue
		}

		logger.Info().Msgf("closed: %s", description)
	}
}
