package domain

import "errors"

var (
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrHashingPassword    = errors.New("error hashing password")
	ErrDBInternal         = errors.New("database internal error")
	ErrInvalidInput       = errors.New("invalid input")
	ErrLoginFailed        = errors.New("invalid user or password")
	ErrDeletedAccount     = errors.New("deleted account")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrPasswordsNotMatch  = errors.New("confirm password does not match")
	ErrPasswordNotChanged = errors.New("password must be different")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrNotFound           = errors.New("not found")
	ErrWrongPassword      = errors.New("incorrect password")

	ErrGuestHasReservations = errors.New("user has active reservations")
	ErrHostHasReservations  = errors.New("user has room(s) with active reservations")
	ErrCannotDeleteAdmin    = errors.New("admin accounts cannot be deleted")
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
