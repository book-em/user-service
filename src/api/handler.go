package api

import (
	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(us service.UserService) Handler {
	return Handler{us}
}

func (h *Handler) RegisterUser(ctx *gin.Context) {

	var dto domain.UserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidValue, err))
		return
	}

	user, err := h.userService.Register(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})

}
