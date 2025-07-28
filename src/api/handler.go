package api

import (
	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(us service.Service) Handler {
	return Handler{us}
}

func (h *Handler) registerUser(ctx *gin.Context) {
	var dto domain.UserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidValue, err))
		return
	}

	user, err := h.service.Register(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) login(ctx *gin.Context) {
	var dto domain.LoginDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidValue, err))
		return
	}

	jwt, err := h.service.Login(dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
