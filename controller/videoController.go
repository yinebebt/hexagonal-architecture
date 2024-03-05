/*package controller will manage route handling for each services*/

package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/Yinebeb-01/simpleapi/entity"
	"gitlab.com/Yinebeb-01/simpleapi/services"
	"gitlab.com/Yinebeb-01/simpleapi/validators"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(*gin.Context) error
	ShowAll(*gin.Context)
	Delete(*gin.Context) error
	Update(*gin.Context) error
}

type controller struct {
	service services.VideoService
}

// validate is used fro custom validation
var cool *validator.Validate

// New is a constructor to initialize new VideoController
func New(service services.VideoService) VideoController {
	cool = validator.New()
	cool.RegisterValidation("is-cool", validators.CoolTitleValidator)
	return &controller{
		service: service,
	}
}

// FindAll will return videos, (use in a GET request endpoint).
func (c *controller) FindAll() []entity.Video {
	return c.service.FindAll()
}

// Save will bind a video-from a POST request body and append it to Videos.
func (c *controller) Save(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		return err
	}

	//validate your custom validator here
	err = cool.Struct(video)
	if err != nil {
		return err
	}
	c.service.Save(video)
	return nil
}

// ShowAll shows the list of videos via some rendered html/css-format
func (c controller) ShowAll(ctx *gin.Context) {
	videos := c.service.FindAll()
	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
		"msg":    "BY Pragmatic review-yina",
	}
	ctx.HTML(http.StatusOK, "index.html", data)
}

func (c *controller) Update(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	video.ID = id

	//validate your custom validator here
	err = cool.Struct(video)
	if err != nil {
		return err
	}
	c.service.Update(video)
	return nil
}

func (c *controller) Delete(ctx *gin.Context) error {
	var video entity.Video
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	video.ID = id
	c.service.Delete(video)
	return nil
}
