package errs

import (
	"fmt"
	"net/http"
)

func NewUserEmailExistsError(email string) *HTTPError {
	return &HTTPError{ErrorCode: EmailAlreadyExists, StatusCode: http.StatusConflict, Message: fmt.Sprintf("User email %s already exists", email)}
}
