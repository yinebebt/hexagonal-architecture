package rest

import (
	"github.com/gin-gonic/gin"
	"log"
)

// CastContext return *gin.Context from an interface.
func CastContext(ctx interface{}) *gin.Context {
	c, ok := ctx.(*gin.Context)
	if !ok {
		log.Printf("unable to assert interface as *gin.Context, got %T", ctx)
		return nil
	}
	return c
}
