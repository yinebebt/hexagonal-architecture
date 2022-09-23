/*package controller will manage route handling for each services*/

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
	"gitlab.com/Yinebeb-01/simpleAPI/services"
	"gitlab.com/Yinebeb-01/simpleAPI/validators"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(*gin.Context) error
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

// FIndAll will return videos, (use in a GET request endpoint).
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
