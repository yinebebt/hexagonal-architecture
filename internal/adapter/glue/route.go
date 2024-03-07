package glue

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func RegisterRoutes(group *gin.RouterGroup, routes []Router, middleware []gin.HandlerFunc) {
	for _, route := range routes {
		var handler []gin.HandlerFunc
		handler = append(handler, middleware...)
		handler = append(handler, route.Handler)
		group.Handle(route.Method, route.Path, handler...)
	}
}
