package httpTransport

import (
	"go-api-template/internal/libs/renderer"
	"go-api-template/internal/service"
	usersHttpTransport "go-api-template/internal/transport/http/users"
)

type HTTPTransport struct {
	Users *usersHttpTransport.UserHandlers
	// others, ex: Orders *ordersHttpTransport.OrderHandlers
}

func NewHTTPTransport(services *service.Services, responseRenderer *renderer.ResponseRenderer) *HTTPTransport {

	return &HTTPTransport{
		Users: usersHttpTransport.NewUserHandlers(services.UserService, responseRenderer),
	}
}
