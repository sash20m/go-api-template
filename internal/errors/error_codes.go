package errs

type ErrorCode int

const (
	Unknown            ErrorCode = 0
	EmailAlreadyExists ErrorCode = 1
)
