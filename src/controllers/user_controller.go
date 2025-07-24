package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"bookem-user-service/utils"

	models "bookem-user-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var payload *models.UserDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	user := models.User{
		Username: payload.Username,
		Password: hashedPassword,
		Email:    strings.ToLower(payload.Email),
		Name:     payload.Name,
		Surname:  payload.Surname,
		Role:     models.UserRole(strings.ToLower(payload.Role)),
		Address:  payload.Address,
	}

	result := uc.DB.Create(&user)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		var msg string
		if strings.Contains(result.Error.Error(), "username") {
			msg = "Username is already taken"
		} else {
			msg = "Email is already taken"
		}
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": msg})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "user": user})
}
