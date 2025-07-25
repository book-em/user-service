package api

import (
	"bookem-user-service/api/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	userHandler UserHandler
}

func NewUserRoute(userHandler UserHandler) *UserRoute {
	return &UserRoute{userHandler}
}

func (ur *UserRoute) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")

	router.Use(middleware.ErrorHandlingMiddleware())
	router.POST("/register", ur.userHandler.RegisterUser)

}
