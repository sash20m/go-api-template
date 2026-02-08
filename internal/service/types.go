package service

import (
	"go-api-template/internal/libs/queue"
	"go-api-template/internal/repositories"

	"github.com/jmoiron/sqlx"
)

type Services struct {
	UserService *UserService
}

func NewServices(db *sqlx.DB, publisher queue.Publisher) *Services {
	userRepository := repositories.NewUserRepository(db)

	userService := NewUserService(userRepository, publisher)

	return &Services{
		UserService: userService,
	}
}

type CreateUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
