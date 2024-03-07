package glue

import (
	"github.com/gin-gonic/gin"
	middlewares2 "gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/middleware"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"net/http"
)

func InitVideoRoute(grp *gin.RouterGroup, video port.VideoHandler) {
	videoRoutes := []Router{
		{
			Method:  http.MethodPost,
			Path:    "/videos",
			Handler: video.Save,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/videos/:id",
			Handler: video.Delete,
		},
		{
			Method:  http.MethodPut,
			Path:    "/videos/:id",
			Handler: video.Update,
		},
		{
			Method: http.MethodGet,
			Path:   "/test",
			Handler: func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"Message": "hello",
				})
			},
		},
	}

	//apiRoute group used to group 'api/*' endpoints.
	RegisterRoutes(grp.Group("/api"), videoRoutes, []gin.HandlerFunc{middlewares2.AuthorizeJWT()})
}
