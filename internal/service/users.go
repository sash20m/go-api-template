package service

import (
	"context"
	"encoding/json"
	"fmt"
	errs "go-api-template/internal/errors"
	"go-api-template/internal/libs/crypto"
	"go-api-template/internal/libs/queue"
	"go-api-template/internal/model"
	"go-api-template/internal/repositories"
)

type UserService struct {
	UserRepository *repositories.UserRepository
	Publisher      queue.Publisher
}

func NewUserService(userRepository *repositories.UserRepository, publisher queue.Publisher) *UserService {
	return &UserService{UserRepository: userRepository, Publisher: publisher}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (any, error) {
	fmt.Println("Getting user by ID: ", id)
	user, err := s.UserRepository.GetUser(ctx, id)

	fmt.Println("User: ", err)
	if err != nil {

		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user CreateUserInput) (model.User, error) {
	emailExists, err := s.UserRepository.CheckIfUserEmailExists(ctx, user.Email)
	if err != nil {
		return model.User{}, err
	}
	if emailExists {
		return model.User{}, errs.NewUserEmailExistsError(user.Email)
	}

	hashedPassword, err := crypto.HashPassword(user.Password)
	if err != nil {
		return model.User{}, err
	}

	userModel := model.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	}
	createdUser, err := s.UserRepository.CreateUser(ctx, userModel)
	if err != nil {
		return model.User{}, err
	}

	if body, err := json.Marshal(map[string]any{
		"id":        createdUser.ID,
		"email":     createdUser.Email,
		"firstName": createdUser.FirstName,
		"lastName":  createdUser.LastName,
	}); err == nil {
		_ = s.Publisher.Publish(ctx, queue.Message{
			Exchange:    queue.EventsExchangeName,
			RoutingKey:  "users.created",
			Body:        body,
			ContentType: "application/json",
		})
	}

	return createdUser, nil
}
