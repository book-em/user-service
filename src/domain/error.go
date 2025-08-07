package domain

import "errors"

var (
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrHashingPassword    = errors.New("error hashing password")
	ErrDBInternal         = errors.New("database internal error")
	ErrInvalidInput       = errors.New("invalid input")
	ErrLoginFailed        = errors.New("invalid user or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrPasswordsNotMatch  = errors.New("confirm password does not match")
	ErrPasswordNotChanged = errors.New("password must be different")
	ErrUnauthenticated    = errors.New("unauthenticated")
)

type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(msg string, code int) error {
	return &AppError{Message: msg, StatusCode: code}
}
