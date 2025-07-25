package api

import (
	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) UserHandler {
	return UserHandler{s}
}

func (h *UserHandler) RegisterUser(ctx *gin.Context) {

	var dto domain.UserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(domain.NewAppError("Invalid input", http.StatusBadRequest))
		return
	}

	user, err := h.service.Register(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})

}
