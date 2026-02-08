package usersHttpTransport

import (
	"go-api-template/internal/transport/http/middlewares"

	"github.com/go-chi/chi/v5"
)

func (h *UserHandlers) RegisterRoutes(r chi.Router) {
	r.Get("/{id}", h.GetUser)
	r.Post("/", middlewares.ChainMiddlewares(h.CreateUser, middlewares.AuthMiddleware))
}
