package api

import (
	"bookem-user-service/api/middleware"
	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	utils "bookem-user-service/util"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(us service.Service) Handler {
	return Handler{us}
}

func (h *Handler) registerUser(ctx *gin.Context) {
	_, span := utils.NewSpan(ctx, "register-user")
	defer span.End()

	var dto domain.UserCreateDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	user, err := h.service.Register(&dto)
	if err != nil {
		ctx.Error(err)
		utils.AddEvent(span, "failed registering user", err)
		return
	}

	ctx.JSON(http.StatusCreated, domain.NewUserDTO(user))
}

func (h *Handler) login(ctx *gin.Context) {
	_, span := utils.NewSpan(ctx, "login-user")
	defer span.End()

	var dto domain.LoginDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	jwt, err := h.service.Login(dto)
	if err != nil {
		ctx.Error(err)
		utils.AddEvent(span, "failed logging in user", err)
		return
	}

	ctx.JSON(http.StatusOK, domain.JWTDTO{Jwt: jwt})
}

func (h *Handler) update(ctx *gin.Context) {
	_, span := utils.NewSpan(ctx, "update-user")
	defer span.End()

	jwt, err := middleware.GetJwt(ctx)
	if err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.AddEvent(span, "unauthenticated", err)
		return
	}

	var dto domain.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	utils.SetSpanUser(span, jwt.ID)

	_, err = h.service.Update(jwt.ID, dto)
	if err != nil {
		ctx.Error(err)
		utils.AddEvent(span, "failed updating user", err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) changePassword(ctx *gin.Context) {
	_, span := utils.NewSpan(ctx, "change-user-password")
	defer span.End()

	jwt, err := middleware.GetJwt(ctx)
	if err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.AddEvent(span, "unauthenticated", err)
		return
	}

	var dto domain.PasswordUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	utils.SetSpanUser(span, jwt.ID)

	_, err = h.service.ChangePassword(jwt.ID, dto)
	if err != nil {
		ctx.Error(err)
		utils.AddEvent(span, "failed changing password", err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) findById(ctx *gin.Context) {
	_, span := utils.NewSpan(ctx, "find-user-by-id")
	defer span.End()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		ctx.Error(err)
		utils.AddEvent(span, "failed parsing ID", err)
		return
	}

	utils.AddAttribInt(span, "id", id)

	log.Printf("Find user by id %d", id)

	user, err := h.service.FindById(uint(id))
	if err != nil {
		ctx.Error(err)
		utils.AddEvent(span, "failed finding user by ID", err)
		return
	}

	ctx.JSON(http.StatusOK, domain.NewUserDTO(user))
}

func (h *Handler) deleteById(ctx *gin.Context) {
	_, span := utils.NewSpan(ctx, "update-user")
	defer span.End()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		ctx.Error(err)
		utils.AddEvent(span, "failed parsing ID", err)
		return
	}

	utils.AddAttribInt(span, "id", id)

	jwt, err := middleware.GetJwt(ctx)
	if err != nil {
		ctx.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.AddEvent(span, "unauthenticatetd", err)
		return
	}

	utils.SetSpanUser(span, jwt.ID)

	err = h.service.Delete(jwt.ID, uint(id))
	if err != nil {
		ctx.Error(err)
		utils.AddEvent(span, "could not delete user", err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
