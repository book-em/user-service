package api

import (
	"bookem-user-service/api/middleware"

	"github.com/gin-gonic/gin"
)

type Route struct {
	handler Handler
}

func NewRoute(handler Handler) *Route {
	return &Route{handler}
}

func (r *Route) Route(rg *gin.RouterGroup) {
	rg.Use(middleware.ErrorHandlingMiddleware())
	rg.POST("/register", r.handler.registerUser)
	rg.POST("/login", r.handler.login)
	rg.PUT("/update", r.handler.update)
	rg.PUT("/password", r.handler.changePassword)
	rg.GET("/:id", r.handler.findById)
	rg.DELETE("/", r.handler.delete)
}
