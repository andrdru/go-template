package middlewares

import (
	"github.com/julienschmidt/httprouter"
)

type HTTPMiddleware func(next httprouter.Handle) httprouter.Handle

func HTTPRouterChain(
	h httprouter.Handle,
	middlewares ...HTTPMiddleware,
) httprouter.Handle {
	if len(middlewares) == 0 {
		return h
	}

	var wrapped = h

	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}

	return wrapped
}
