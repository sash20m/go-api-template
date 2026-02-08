package renderer

import (
	"errors"
	errs "go-api-template/internal/errors"
	"net/http"
	"time"

	"github.com/unrolled/render"
)

// SuccessfulResponse standardizes responses with 200-299 status code
type SuccessfulResponse struct {
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// ClientErrorResponse standardizes responses with 400-499 status code
type ClientErrorResponse struct {
	errs.HTTPError
	Timestamp time.Time `json:"timestamp"`
}

// ServerErrorResponse standardizes responses with 500-599 status code
type ServerErrorResponse struct {
	errs.HTTPError
	Timestamp time.Time `json:"timestamp"`
}

type ResponseRenderer struct {
	Render *render.Render
}

func NewResponseRenderer() *ResponseRenderer {
	return &ResponseRenderer{Render: render.New(render.Options{IndentJSON: true})}
}

func (s *ResponseRenderer) JSON(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")

	if httpErr, ok := asHTTPError(response); ok {
		s.JSONCustomHTTPError(w, httpErr)
		return
	}

	if err, ok := response.(error); ok && err != nil {
		s.JSONInternalError(w)
		return
	}

	if statusCode/100 == 2 {
		response := SuccessfulResponse{Data: response, Timestamp: time.Now().UTC()}
		err := s.Render.JSON(w, statusCode, response)
		if err != nil {
			panic(err)
		}

		return
	}

	err := s.Render.JSON(w, statusCode, response)
	if err != nil {
		panic(err)
	}
}

// Other formats... (xml, html etc)

// JSONCustomHTTPError renders errors produced by the app (errs.HTTPError) in a consistent envelope.
func (s *ResponseRenderer) JSONCustomHTTPError(w http.ResponseWriter, httpErr *errs.HTTPError) {
	if httpErr == nil {
		s.JSONInternalError(w)
		return
	}

	statusCode := httpErr.StatusCode
	codeClass := statusCode / 100

	if codeClass == 4 {
		response := ClientErrorResponse{HTTPError: *httpErr, Timestamp: time.Now().UTC()}
		if err := s.Render.JSON(w, statusCode, response); err != nil {
			panic(err)
		}
		return
	}

	// Default to server error for anything else (including malformed status codes).
	s.JSONInternalError(w)
}

func (s *ResponseRenderer) JSONInternalError(w http.ResponseWriter) {
	httpError := errs.HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error.",
	}

	response := ServerErrorResponse{HTTPError: httpError, Timestamp: time.Now().UTC()}

	if err := s.Render.JSON(w, httpError.StatusCode, response); err != nil {
		panic(err)
	}
}

func asHTTPError(v any) (*errs.HTTPError, bool) {
	if v == nil {
		return nil, false
	}

	// If it's already *HTTPError
	if he, ok := v.(*errs.HTTPError); ok && he != nil {
		return he, true
	}

	// If it's an error (possibly wrapped), unwrap to *HTTPError
	if e, ok := v.(error); ok && e != nil {
		var he *errs.HTTPError
		if errors.As(e, &he) && he != nil {
			return he, true
		}
	}

	return nil, false
}
