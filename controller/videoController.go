/*package controller will manage route handling for each services*/

package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
	"gitlab.com/Yinebeb-01/simpleAPI/services"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(*gin.Context) entity.Video
}

type controller struct {
	service services.VideoService
}

//New is a constructor to initialize new VideoController
func New(service services.VideoService) VideoController {
	return &controller{
		service: service,
	}
}

//FIndAll will return videos, (use in a GET request endpoint).
func (c *controller) FindAll() []entity.Video {
	return c.service.FindAll()
}

//Save will bind a video-from a POST request body and append it to Videos.
func (c *controller) Save(ctx *gin.Context) entity.Video {
	var video entity.Video
	ctx.BindJSON(&video)
	c.service.Save(video)
	return video
}
