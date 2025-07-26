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
	rg.POST("/register", r.handler.RegisterUser)
}
