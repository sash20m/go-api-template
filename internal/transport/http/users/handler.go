package usersHttpTransport

import (
	"go-api-template/internal/libs/renderer"
	"go-api-template/internal/libs/utils"
	"go-api-template/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandlers struct {
	userService      *service.UserService
	responseRenderer *renderer.ResponseRenderer
}

func NewUserHandlers(userService *service.UserService, responseRenderer *renderer.ResponseRenderer) *UserHandlers {
	return &UserHandlers{userService: userService, responseRenderer: responseRenderer}
}

func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		h.responseRenderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	h.responseRenderer.JSON(w, http.StatusOK, user)
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetJSONBody[CreateUserRequest](r)
	if err != nil {
		h.responseRenderer.JSON(w, http.StatusBadRequest, err)
		return
	}

	userInput := toCreateUserInput(user)
	createdUser, err := h.userService.CreateUser(r.Context(), userInput)
	if err != nil {
		h.responseRenderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	h.responseRenderer.JSON(w, http.StatusOK, createdUser)
}
