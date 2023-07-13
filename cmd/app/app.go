package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/andrdru/go-template/configs"
	"github.com/andrdru/go-template/internal/api"
	"github.com/andrdru/go-template/internal/managers"
	"github.com/andrdru/go-template/internal/repos"
)

type (
	//closer resource graceful shutdown handler
	closer func(ctx context.Context) (description string, err error)

	bootstrap struct {
		httpListenAndServe func()

		closers []closer
	}
)

func Run(logger zerolog.Logger, configPath string) (code int) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info().Msg("app staring")
	defer logger.Info().Msg("app finished")

	conf, err := configs.NewConfig(configPath)
	if err != nil {
		logger.Error().Err(err).Msg("init config")
		return 1
	}

	boot, err := initApp(logger, conf)
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

	go boot.httpListenAndServe()

	logger.Info().Msg("app started successfully")
	<-ctx.Done()

	return 0
}

func initApp(logger zerolog.Logger, conf configs.Config) (boot bootstrap, err error) {
	boot = bootstrap{}

	db, err := conf.Postgres.Connect()
	if err != nil {
		return bootstrap{}, fmt.Errorf("postgres connect: %w", err)
	}

	userRepo := repos.NewUser(db)
	authManager := managers.NewAuth(userRepo)

	httpAPI := api.NewAPI(logger, authManager)
	router := httpAPI.InitRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.HTTP.Host, conf.HTTP.Port),
		Handler: router,
	}

	boot.httpListenAndServe = func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("serve http")
		}
	}

	boot.closers = append(boot.closers, func(ctx context.Context) (description string, err error) {
		err = srv.Shutdown(ctx)
		return "http server", err
	})

	return boot, nil
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
