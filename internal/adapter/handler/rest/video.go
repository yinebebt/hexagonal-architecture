package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"github.com/gin-gonic/gin"
)

type video struct {
	videoService service.VideoService
}

// InitVideo is a constructor to initialize VideoHandler
func InitVideo(videoSer service.VideoService) port.VideoHandler {
	return &video{
		videoService: videoSer,
	}
}

// Save will bind a video-from a POST request body and append it to Videos.
//
//	@Summary		Save video
//	@Description	Save video description
//	@Tags			Video
//	@Accept			json
//	@Param			video	body	entity.Video	true	"video to save"
//	@Produce		json
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos [post]
func (v *video) Save(ctx interface{}) {
	ginCtx := CastContext(ctx)
	req := entity.Video{}
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		err = errors.New("invalid input: " + err.Error())
		ginCtx.JSON(http.StatusBadRequest, err.Error())
		_ = ginCtx.Error(err)
		return
	}

	vid, err := v.videoService.Save(req)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, gin.H{"data": vid})
}

// FindAll will return videos, (use in a GET request endpoint).
//
//	@Summary		FindAll video
//	@Description	FindAll video description
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos [get]
func (v *video) FindAll(ctx interface{}) {
	ginCtx := CastContext(ctx)
	res, err := v.videoService.FindAll()
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, res)
}

// Delete
//
//	@Summary		Delete video
//	@Description	Delete video description
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos/:id [delete]
func (v *video) Delete(ctx interface{}) {
	ginCtx := CastContext(ctx)
	id, err := strconv.ParseInt(ginCtx.Param("id"), 0, 0)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
	v.videoService.Delete(id)
	ginCtx.JSON(http.StatusOK, gin.H{"message": "video deleted successfully!"})
}

// Update
//
//	@Summary		Update video
//	@Description	Update video description
//	@Tags			Video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Video
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/videos [put]
func (v *video) Update(ctx interface{}) {
	ginCtx := CastContext(ctx)
	var video entity.Video
	err := ginCtx.ShouldBindJSON(&video)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	id, err := strconv.ParseInt(ginCtx.Param("id"), 0, 0)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
	video.ID = id
	res := v.videoService.Update(video)
	ginCtx.JSON(http.StatusCreated, gin.H{"data": res})
}
