package usersHttpTransport

import "go-api-template/internal/service"

func toCreateUserInput(request CreateUserRequest) service.CreateUserInput {
	return service.CreateUserInput{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}
}
