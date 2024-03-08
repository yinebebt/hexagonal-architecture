package route

import (
	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/glue"
	middlewares2 "github.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/middleware"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitVideoRoute(grp *gin.RouterGroup, video port.VideoHandler) {
	videoRoutes := []glue.Router{
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
	viewRoutes := []glue.Router{
		{
			Method:  http.MethodGet,
			Handler: video.ShowAll,
			Path:    "/videos",
		},
	}

	//apiRoute group used to group 'api/*' endpoints.
	glue.RegisterRoutes(grp.Group(""), videoRoutes, []gin.HandlerFunc{middlewares2.AuthorizeJWT()})
	//viewRoute Group will use to render static files
	glue.RegisterRoutes(grp.Group("/view"), viewRoutes, []gin.HandlerFunc{})
}
