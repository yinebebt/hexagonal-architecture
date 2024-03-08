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
			Handler: handler.Login,
		},
	}

	// Login Endpoint: Authentication
	glue.RegisterRoutes(grp.Group(""), loginRoutes, []gin.HandlerFunc{})
}
