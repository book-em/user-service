package api

import (
	"bookem-user-service/api/middleware"
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
	var dto domain.UserCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		return
	}

	user, err := h.service.Register(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	result := domain.UserDTO{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Surname:  user.Surname,
		Address:  user.Address,
		Role:     string(user.Role),
	}

	ctx.JSON(http.StatusCreated, result)
}

func (h *Handler) login(ctx *gin.Context) {
	var dto domain.LoginDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		return
	}

	jwt, err := h.service.Login(dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, domain.JWTDTO{Jwt: jwt})
}

func (h *Handler) update(ctx *gin.Context) {
	jwt, err := middleware.GetJwt(ctx)
	if err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		return
	}

	var dto domain.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		return
	}

	_, err = h.service.Update(jwt.ID, dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) changePassword(ctx *gin.Context) {
	jwt, err := middleware.GetJwt(ctx)
	if err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		return
	}

	var dto domain.PasswordUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		return
	}

	_, err = h.service.ChangePassword(jwt.ID, dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
