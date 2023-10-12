package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrdru/go-template/internal/pkg/graceful"

	"github.com/andrdru/go-template/configs"
	"github.com/andrdru/go-template/internal/api"
	"github.com/andrdru/go-template/internal/managers"
	"github.com/andrdru/go-template/internal/repos"
)

type (
	bootstrap struct {
		httpListenAndServe func()

		closers []graceful.Closer
	}
)

func Run(logger *slog.Logger, configPath string) (code int) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("app staring")
	defer logger.Info("app finished")

	conf, err := configs.NewConfig(configPath)
	if err != nil {
		logger.Error("init config", slog.Any("error", err))
		return 1
	}

	boot, err := initApp(logger, conf)
	if err != nil {
		logger.Error("init app", slog.Any("error", err))
		return 1
	}

	defer func() {
		ctxCloser, cancelCloser := context.WithTimeout(context.Background(), 15*time.Second)
		defer func() {
			cancelCloser()
		}()

		graceful.Stop(ctxCloser, logger, boot.closers)
	}()

	go boot.httpListenAndServe()

	logger.Info("app started successfully")
	<-ctx.Done()

	return 0
}

func initApp(logger *slog.Logger, conf configs.Config) (boot bootstrap, err error) {
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
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("serve http", slog.Any("error", err))
		}
	}

	boot.closers = append(boot.closers, func(ctx context.Context) (description string, err error) {
		err = srv.Shutdown(ctx)
		return "http server", err
	})

	return boot, nil
}
