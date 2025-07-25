package middleware

import (
	domain "bookem-user-service/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUserExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrEmailExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrHashingPassword):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}

		status := MapErrorToStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
	}
}
