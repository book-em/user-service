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
