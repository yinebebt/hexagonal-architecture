package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/simpleAPI/controller"
	"gitlab.com/Yinebeb-01/simpleAPI/services"
)

var (
	videoservice    services.VideoService      = services.New()
	videocontroller controller.VideoController = controller.New(videoservice)
)

func main() {
	router := gin.Default()
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Message": "hello",
		})
	})

	router.GET("/videos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, videocontroller.FindAll())
	})

	router.POST("/videos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, videocontroller.Save(ctx))
	})

	router.Run(":8080")
}
