package api

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/andrdru/go-template/internal/entities"
	"github.com/andrdru/go-template/internal/middlewares"
)

type (
	API struct {
		logger zerolog.Logger

		authManager authManager
	}

	authManager interface {
		Check(r *http.Request) (ctx context.Context, err error)
		Login(ctx context.Context, w http.ResponseWriter, session entities.Session) error
	}
)

var (
	OptInternalError = Error("internal error")
	OptUnauthorized  = Error("unauthorized")
)

func NewAPI(logger zerolog.Logger, sessionManager authManager) *API {
	return &API{
		logger:      logger,
		authManager: sessionManager,
	}
}

func (a *API) InitRoutes() *httprouter.Router {
	router := initHTTP()

	auth := []middlewares.HTTPMiddleware{
		middlewares.SessionValidate(a.logger, a.authManager, handleUnauthorized),
	}

	// anonymous methods
	router.Handle(http.MethodPost, "/user/authorize", a.UserAuthorize)

	// auth methods
	router.Handle(http.MethodGet, "/user/:id", middlewares.HTTPRouterChain(a.UserGet, auth...))

	return router
}

func initHTTP() *httprouter.Router {
	router := httprouter.New()

	router.GET("/health", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		_, _ = w.Write([]byte("ok"))
	})

	router.GET("/metrics", httpHandlerAdapterFunc(promhttp.Handler()))

	router.HandlerFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)
	router.Handler(http.MethodGet, "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handler(http.MethodGet, "/debug/pprof/heap", pprof.Handler("heap"))
	router.Handler(http.MethodGet, "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handler(http.MethodGet, "/debug/pprof/block", pprof.Handler("block"))

	return router
}

// httpHandlerAdapterFunc httprouter adapter for http.Handler
func httpHandlerAdapterFunc(handler http.Handler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler.ServeHTTP(writer, request)
	}
}

func handleUnauthorized(w http.ResponseWriter, message string) error {
	m := NewMessage()
	m.SetError(Code(http.StatusUnauthorized), OptUnauthorized)
	if message != "" {
		m.SetError(Error(message))
	}

	return m.Return(w)
}
