package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type video struct {
	videoService service.VideoService
}

var cool *validator.Validate // custom validation

// InitVideo is a constructor to initialize VideoHandler
func InitVideo(videoSer service.VideoService) port.VideoHandler {
	cool = validator.New()
	err := cool.RegisterValidation("is-cool", util.CoolTitleValidator)
	if err != nil {
		return nil
	}
	return &video{
		videoService: videoSer,
	}
}

// Save will bind a video-from a POST request body and append it to Videos.
// @Summary      Save video
// @Description  Save video description
// @Tags         Video
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.Video
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /videos [post]
func (v *video) Save(ctx interface{}) {
	ginCtx := CastContext(ctx)
	req := entity.Video{}
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		err = errors.New("invalid input")
		_ = ginCtx.Error(err)
		return
	}

	vid, err := v.videoService.Save(req)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ginCtx.JSON(http.StatusOK, gin.H{"data": vid})
	}
}

// FindAll will return videos, (use in a GET request endpoint).
// @Summary      FindAll video
// @Description  FindAll video description
// @Tags         Video
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.Video
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /videos [get]
func (v *video) FindAll(ctx interface{}) {
	ginCtx := CastContext(ctx)
	res, err := v.videoService.FindAll()
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, res)
}

// ShowAll shows the list of videos via some rendered html/css-format
// @Summary      ShowAll video
// @Description  Show all video description
// @Tags         Video
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.Video
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /all-videos [get]
// Deprecated: ShowAll is deprecated, use FindAll instead
func (v *video) ShowAll(ctx interface{}) {
	ginCtx := CastContext(ctx)
	videos, err := v.videoService.FindAll()
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
		"msg":    "By st-son admin",
	}
	ginCtx.HTML(http.StatusOK, "index.html", data)
}

// Delete
// @Summary      Delete video
// @Description  Delete video description
// @Tags         Video
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.Video
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /videos [delete]
func (v *video) Delete(ctx interface{}) {
	ginCtx := CastContext(ctx)
	video := entity.Video{}
	id, err := strconv.ParseUint(ginCtx.Param("id"), 0, 0)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
	video.ID = id
	v.videoService.Delete(video)
	ginCtx.JSON(http.StatusOK, gin.H{"message": "video deleted successfully!"})
}

// Update
// @Summary      Update video
// @Description  Update video description
// @Tags         Video
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.Video
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /videos [put]
func (v *video) Update(ctx interface{}) {
	ginCtx := CastContext(ctx)
	var video entity.Video
	err := ginCtx.ShouldBindJSON(&video)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	id, err := strconv.ParseUint(ginCtx.Param("id"), 0, 0)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
	video.ID = id

	//validate your custom validator here
	err = cool.Struct(video)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}
	res := v.videoService.Update(video)
	ginCtx.JSON(http.StatusCreated, gin.H{"data": res})
}
