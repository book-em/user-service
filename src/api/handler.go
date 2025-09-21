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

func (h *Handler) registerUser(c *gin.Context) {
	ctx, span := utils.NewSpan(c.Request.Context(), "register-user")
	defer span.End()

	var dto domain.UserCreateDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	user, err := h.service.Register(ctx, &dto)
	if err != nil {
		c.Error(err)
		utils.AddEvent(span, "failed registering user", err)
		return
	}

	c.JSON(http.StatusCreated, domain.NewUserDTO(user))
}

func (h *Handler) login(c *gin.Context) {
	_, span := utils.NewSpan(c.Request.Context(), "login-user")
	defer span.End()

	var dto domain.LoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	jwt, err := h.service.Login(dto)
	if err != nil {
		c.Error(err)
		utils.AddEvent(span, "failed logging in user", err)
		return
	}

	c.JSON(http.StatusOK, domain.JWTDTO{Jwt: jwt})
}

func (h *Handler) update(c *gin.Context) {
	_, span := utils.NewSpan(c.Request.Context(), "update-user")
	defer span.End()

	jwt, err := middleware.GetJwt(c)
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.AddEvent(span, "unauthenticated", err)
		return
	}

	var dto domain.UserUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	utils.SetSpanUser(span, jwt.ID)

	_, err = h.service.Update(jwt.ID, dto)
	if err != nil {
		c.Error(err)
		utils.AddEvent(span, "failed updating user", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) changePassword(c *gin.Context) {
	_, span := utils.NewSpan(c.Request.Context(), "change-user-password")
	defer span.End()

	jwt, err := middleware.GetJwt(c)
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.AddEvent(span, "unauthenticated", err)
		return
	}

	var dto domain.PasswordUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.AddEvent(span, "failed binding JSON", err)
		return
	}

	utils.SetSpanUser(span, jwt.ID)

	_, err = h.service.ChangePassword(jwt.ID, dto)
	if err != nil {
		c.Error(err)
		utils.AddEvent(span, "failed changing password", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) findById(c *gin.Context) {
	_, span := utils.NewSpan(c.Request.Context(), "find-user-by-id")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		c.Error(err)
		utils.AddEvent(span, "failed parsing ID", err)
		return
	}

	utils.AddAttribInt(span, "id", id)

	log.Printf("Find user by id %d", id)

	user, err := h.service.FindById(uint(id))
	if err != nil {
		c.Error(err)
		utils.AddEvent(span, "failed finding user by ID", err)
		return
	}

	c.JSON(http.StatusOK, domain.NewUserDTO(user))
}

func (h *Handler) deleteById(c *gin.Context) {
	_, span := utils.NewSpan(c.Request.Context(), "update-user")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		c.Error(err)
		utils.AddEvent(span, "failed parsing ID", err)
		return
	}

	utils.AddAttribInt(span, "id", id)

	jwt, err := middleware.GetJwt(c)
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.AddEvent(span, "unauthenticatetd", err)
		return
	}

	utils.SetSpanUser(span, jwt.ID)

	err = h.service.Delete(jwt.ID, uint(id))
	if err != nil {
		c.Error(err)
		utils.AddEvent(span, "could not delete user", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
