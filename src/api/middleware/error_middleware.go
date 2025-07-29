package middleware

import (
	domain "bookem-user-service/domain"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUsernameExists), errors.Is(err, domain.ErrEmailExists):
		return http.StatusConflict
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

		// Some debug variable or mode can be used
		log.Printf("[DEBUG] Error: %v\n", err)

		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}

		status := mapErrorToStatus(err)
		c.JSON(status, gin.H{"error": err.Error()})
	}
}
