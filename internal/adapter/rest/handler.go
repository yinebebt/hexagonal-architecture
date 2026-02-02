package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
	"github.com/yinebebt/hexagonal-architecture/internal/core/service"
)

// videoHandler implements port.VideoHandler for REST API
type videoHandler struct {
	videoService service.VideoService
}

// NewVideoHandler creates a new REST video handler that implements port.VideoHandler
func NewVideoHandler(videoService service.VideoService) port.VideoHandler {
	return &videoHandler{
		videoService: videoService,
	}
}

// Save handles POST /videos - creates a new video
//
//	@Summary		Save video
//	@Description	Save video description
//	@Tags			Video
//	@Accept			json
//	@Param			video	body	entity.Video	true	"video to save"
//	@Produce		json
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos [post]
func (v *videoHandler) Save(ctx interface{}) {
	ginCtx := castContext(ctx)
	req := entity.Video{}
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	vid, err := v.videoService.Save(req)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, gin.H{"data": vid})
}

// FindAll handles GET /videos - returns all videos
//
//	@Summary		FindAll video
//	@Description	FindAll video description
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Video
//	@Failure		500	{object}	map[string]string
//	@Router			/videos [get]
func (v *videoHandler) FindAll(ctx interface{}) {
	ginCtx := castContext(ctx)
	res, err := v.videoService.FindAll()
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, res)
}

// Delete handles DELETE /videos/:id - deletes a video
//
//	@Summary		Delete video
//	@Description	Delete video description
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos/:id [delete]
func (v *videoHandler) Delete(ctx interface{}) {
	ginCtx := castContext(ctx)
	id, err := strconv.ParseInt(ginCtx.Param("id"), 10, 64)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	err = v.videoService.Delete(id)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, gin.H{"message": "video deleted successfully"})
}

// Update handles PUT /videos/:id - updates a video
//
//	@Summary		Update video
//	@Description	Update video description
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Video ID"
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos/:id [put]
func (v *videoHandler) Update(ctx interface{}) {
	ginCtx := castContext(ctx)
	var video entity.Video
	err := ginCtx.ShouldBindJSON(&video)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	id, err := strconv.ParseInt(ginCtx.Param("id"), 10, 64)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}
	video.ID = id
	res, err := v.videoService.Update(video)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, gin.H{"data": res})
}

// castContext returns *gin.Context from an interface.
// It panics if the context cannot be cast, as this indicates a programming error.
func castContext(ctx interface{}) *gin.Context {
	c, ok := ctx.(*gin.Context)
	if !ok {
		panic("unable to assert interface as *gin.Context")
	}
	return c
}

// Route represents a single HTTP route definition
type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// RegisterRoutes registers multiple routes to a Gin router group
func RegisterRoutes(group *gin.RouterGroup, routes []Route) {
	for _, route := range routes {
		group.Handle(route.Method, route.Path, route.Handler)
	}
}

// RegisterVideoRoutes registers video-related REST routes
func RegisterVideoRoutes(grp *gin.RouterGroup, handler port.VideoHandler) {
	routes := []Route{
		{
			Method:  http.MethodPost,
			Path:    "/videos",
			Handler: handlerToGinFunc(handler, "save"),
		},
		{
			Method:  http.MethodGet,
			Path:    "/videos",
			Handler: handlerToGinFunc(handler, "find_all"),
		},
		{
			Method:  http.MethodPut,
			Path:    "/videos/:id",
			Handler: handlerToGinFunc(handler, "update"),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/videos/:id",
			Handler: handlerToGinFunc(handler, "delete"),
		},
	}

	RegisterRoutes(grp, routes)
}

// handlerToGinFunc converts port.VideoHandler to gin.HandlerFunc
func handlerToGinFunc(handler port.VideoHandler, action string) gin.HandlerFunc {
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
