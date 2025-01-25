package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yinebebt/hexagonal-architecture/internal/adapter/glue"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
	"net/http"
)

func InitVideoRoute(grp *gin.RouterGroup, video port.VideoHandler) {
	videoRoutes := []glue.Router{
		{
			Method:  http.MethodPost,
			Path:    "/videos",
			Handler: VideoHandlerFunc(video, "save"),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/videos/:id",
			Handler: VideoHandlerFunc(video, "delete"),
		},
		{
			Method:  http.MethodPut,
			Path:    "/videos/:id",
			Handler: VideoHandlerFunc(video, "update"),
		},
		{
			Method:  http.MethodGet,
			Handler: VideoHandlerFunc(video, "find_all"),
			Path:    "/videos",
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
	glue.RegisterRoutes(grp.Group(""), videoRoutes)
}

// VideoHandlerFunc converts VideoHandler adapter to gin.HandlerFunc
func VideoHandlerFunc(handler port.VideoHandler, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch action {
		case "save":
			handler.Save(c)
		case "find_all":
			handler.FindAll(c)
		case "update":
			handler.Update(c)
		case "delete":
			handler.Delete(c)

		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid action"})
		}
	}
}
