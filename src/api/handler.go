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
	"go.opentelemetry.io/otel/attribute"
)

type Handler struct {
	service service.Service
}

func NewHandler(us service.Service) Handler {
	return Handler{us}
}

func (h *Handler) registerUser(c *gin.Context) {
	utils.TEL.Push(c.Request.Context(), "register-user")
	defer utils.TEL.Pop()

	var dto domain.UserCreateDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.TEL.Event("failed binding JSON", err)
		return
	}

	user, err := h.service.Register(utils.TEL.Ctx(), &dto)
	if err != nil {
		c.Error(err)
		utils.TEL.Event("failed registering user", err)
		return
	}

	c.JSON(http.StatusCreated, domain.NewUserDTO(user))
}

func (h *Handler) login(c *gin.Context) {
	utils.TEL.Push(c.Request.Context(), "login-user")
	defer utils.TEL.Pop()

	var dto domain.LoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.TEL.Event("failed binding JSON", err)
		return
	}

	jwt, err := h.service.Login(utils.TEL.Ctx(), dto)
	if err != nil {
		c.Error(err)
		utils.TEL.Event("failed logging in user", err)
		return
	}

	c.JSON(http.StatusOK, domain.JWTDTO{Jwt: jwt})
}

func (h *Handler) update(c *gin.Context) {
	utils.TEL.Push(c.Request.Context(), "update-user")
	defer utils.TEL.Pop()

	jwt, err := middleware.GetJwt(c)
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.TEL.Event("unauthenticated", err)
		return
	}

	var dto domain.UserUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.TEL.Event("failed binding JSON", err)
		return
	}

	utils.TEL.SetUser(jwt.ID)

	_, err = h.service.Update(utils.TEL.Ctx(), jwt.ID, dto)
	if err != nil {
		c.Error(err)
		utils.TEL.Event("failed updating user", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) changePassword(c *gin.Context) {
	utils.TEL.Push(c.Request.Context(), "change-user-password")
	defer utils.TEL.Pop()

	jwt, err := middleware.GetJwt(c)
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.TEL.Event("unauthenticated", err)
		return
	}

	var dto domain.PasswordUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrInvalidInput, err))
		utils.TEL.Event("failed binding JSON", err)
		return
	}

	utils.TEL.SetUser(jwt.ID)

	_, err = h.service.ChangePassword(utils.TEL.Ctx(), jwt.ID, dto)
	if err != nil {
		c.Error(err)
		utils.TEL.Event("failed changing password", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) findById(c *gin.Context) {
	utils.TEL.Push(c.Request.Context(), "find-user-by-id")
	defer utils.TEL.Pop()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		c.Error(err)
		utils.TEL.Event("failed parsing ID", err)
		return
	}

	utils.TEL.SetAttrib(attribute.Int("id", id))

	log.Printf("Find user by id %d", id)

	user, err := h.service.FindById(utils.TEL.Ctx(), uint(id))
	if err != nil {
		c.Error(err)
		utils.TEL.Event("failed finding user by ID", err)
		return
	}

	c.JSON(http.StatusOK, domain.NewUserDTO(user))
}

func (h *Handler) deleteById(c *gin.Context) {
	utils.TEL.Push(c.Request.Context(), "update-user")
	defer utils.TEL.Pop()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		c.Error(err)
		utils.TEL.Event("failed parsing ID", err)
		return
	}

	utils.TEL.SetAttrib(attribute.Int("id", id))

	jwt, err := middleware.GetJwt(c)
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", domain.ErrUnauthenticated, err))
		utils.TEL.Event("unauthenticatetd", err)
		return
	}

	utils.TEL.SetUser(jwt.ID)

	err = h.service.Delete(utils.TEL.Ctx(), jwt.ID, uint(id))
	if err != nil {
		c.Error(err)
		utils.TEL.Event("could not delete user", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
