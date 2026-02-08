package errs

type HTTPError struct {
	ErrorCode  ErrorCode `json:"errorCode,omitempty"`
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(errorCode ErrorCode, statusCode int, msg string) error {
	return &HTTPError{
		ErrorCode:  errorCode,
		StatusCode: statusCode,
		Message:    msg,
	}
}
