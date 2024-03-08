package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/util"
	"net/http"
	"strconv"
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
func (v *video) Save(ctx *gin.Context) {
	req := entity.Video{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = errors.New("invalid input")
		_ = ctx.Error(err)
		return
	}

	vid, err := v.videoService.Save(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": vid})
	}
}

// FindAll will return videos, (use in a GET request endpoint).
func (v *video) FindAll(ctx *gin.Context) {
	res := v.videoService.FindAll()
	ctx.JSON(http.StatusOK, res)
}

// ShowAll shows the list of videos via some rendered html/css-format
func (v *video) ShowAll(ctx *gin.Context) {
	videos := v.videoService.FindAll()
	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
		"msg":    "By st-son admin",
	}
	ctx.HTML(http.StatusOK, "index.html", data)
}

func (v *video) Delete(ctx *gin.Context) {
	video := entity.Video{}
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	video.ID = id
	v.videoService.Delete(video)
	ctx.JSON(http.StatusOK, gin.H{"message": "video deleted successfully!"})
	return
}

func (v *video) Update(ctx *gin.Context) {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	video.ID = id

	//validate your custom validator here
	err = cool.Struct(video)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	res := v.videoService.Update(video)
	ctx.JSON(http.StatusCreated, gin.H{"data": res})
}
