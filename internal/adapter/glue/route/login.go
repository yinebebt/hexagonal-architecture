package route

import (
	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/glue"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitLoginRoute(grp *gin.RouterGroup, handler port.LoginHandler) {
	loginRoutes := []glue.Router{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: LoginHandlerFunc(handler, "login"),
		},
	}

	// Login Endpoint: Authentication
	glue.RegisterRoutes(grp.Group(""), loginRoutes, []gin.HandlerFunc{})
}

// LoginHandlerFunc converts LoginHandler adapter to gin.HandlerFunc
func LoginHandlerFunc(handler port.LoginHandler, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch action {
		case "login":
			handler.Login(c)

		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid action"})
		}
	}
}
