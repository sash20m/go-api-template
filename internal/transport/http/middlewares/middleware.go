package middlewares

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddlewares(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
