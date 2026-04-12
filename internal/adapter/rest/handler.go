package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

// videoHandler implements REST API for video operations
type videoHandler struct {
	videoService port.VideoService
}

// NewVideoHandler creates a new REST video handler
func NewVideoHandler(videoService port.VideoService) *videoHandler {
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
func (v *videoHandler) Save(c *gin.Context) {
	req := entity.Video{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	vid, err := v.videoService.Save(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vid})
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
func (v *videoHandler) FindAll(c *gin.Context) {
	res, err := v.videoService.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// FindByID handles GET /videos/:id - returns a single video
//
//	@Summary		FindByID video
//	@Description	Find a video by its ID
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Video ID"
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/videos/:id [get]
func (v *videoHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	video, err := v.videoService.FindByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": video})
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
func (v *videoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	err = v.videoService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "video deleted successfully"})
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
func (v *videoHandler) Update(c *gin.Context) {
	var video entity.Video
	err := c.ShouldBindJSON(&video)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}
	video.ID = id
	res, err := v.videoService.Update(c.Request.Context(), video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// RegisterVideoRoutes registers video-related REST routes
func RegisterVideoRoutes(grp *gin.RouterGroup, handler *videoHandler) {
	grp.POST("/videos", handler.Save)
	grp.GET("/videos", handler.FindAll)
	grp.GET("/videos/:id", handler.FindByID)
	grp.PUT("/videos/:id", handler.Update)
	grp.DELETE("/videos/:id", handler.Delete)
}
