package domain

import "errors"

var (
	ErrUsernameExists  = errors.New("Username already exists")
	ErrEmailExists     = errors.New("Email already exists")
	ErrHashingPassword = errors.New("Error hashing password")
	ErrDBInternal      = errors.New("Database internal error")
	ErrInvalidInput    = errors.New("Invalid input")
	ErrLoginFailed     = errors.New("Invalid user or password")
	ErrUnauthorized    = errors.New("Unauthorized")
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
