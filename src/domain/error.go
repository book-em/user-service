package domain

import "errors"

var (
	ErrInvalidEmail    = errors.New("Invalid email format")
	ErrUserExists      = errors.New("Username already exists")
	ErrEmailExists     = errors.New("Email already exists")
	ErrHashingPassword = errors.New("Error hashing password")
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
