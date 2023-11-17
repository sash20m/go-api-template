package httputils

import (
	"fmt"
	"net/http"
	"time"

	localErrs "github.com/sash20m/go-api-template/internal/errors"
	"github.com/unrolled/render"
)

// SuccessfulResponse standardizes responses with 200-299 status code
type SuccessfulResponse struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// ClientErrorResponse standardizes responses with 400-499 status code
type ClientErrorResponse struct {
	StatusCode int           `json:"statusCode"`
	Error      localErrs.Err `json:"error"`
	Timestamp  time.Time     `json:"timestamp"`
}

// ServerErrorResponse standardizes responses with 500-599 status code
type ServerErrorResponse struct {
	StatusCode int           `json:"statusCode"`
	Error      localErrs.Err `json:"error"`
	Timestamp  time.Time     `json:"timestamp"`
}

type Sender struct {
	Render *render.Render
}

// JSON formats the v in a json format and sends it to the client with w. If the statusCode
// is specific for a error, then v is assumed to be a string, an error or an custom errors.Err struct
func (s *Sender) JSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	codeClass := statusCode / 100

	if codeClass == 2 {
		response := SuccessfulResponse{Data: v, Timestamp: time.Now().UTC()}
		err := s.Render.JSON(w, statusCode, response)
		return err
	}

	if codeClass == 4 {
		errResponse := localErrs.Err{Message: fmt.Sprint(v)}

		if errInfo, ok := v.(localErrs.Err); ok {
			errResponse.Message = errInfo.Error()
			errResponse.Data = errInfo.Data
		}

		response := ClientErrorResponse{StatusCode: statusCode, Error: errResponse, Timestamp: time.Now().UTC()}
		err := s.Render.JSON(w, statusCode, response)
		return err
	}

	if codeClass == 5 {
		errResponse := localErrs.Err{Message: fmt.Sprint(v)}

		if errInfo, ok := v.(localErrs.Err); ok {
			errResponse.Message = errInfo.Error()
			errResponse.Data = errInfo.Data
		}

		response := ServerErrorResponse{StatusCode: statusCode, Error: errResponse, Timestamp: time.Now().UTC()}
		err := s.Render.JSON(w, statusCode, response)
		return err
	}

	err := s.Render.JSON(w, statusCode, v)
	if err != nil {
		return err
	}

	return nil
}

// Other formats... (xml, html etc)
