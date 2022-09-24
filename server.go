package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/simpleAPI/controller"
	"gitlab.com/Yinebeb-01/simpleAPI/middlewares"
	"gitlab.com/Yinebeb-01/simpleAPI/services"

	dump "github.com/tpkeeper/gin-dump"
)

var (
	videoservice    services.VideoService      = services.New()
	videocontroller controller.VideoController = controller.New(videoservice)
)

func main() {
	configOutput()

	router := gin.New()
	//custome middlewares used, dump is an alise of gin-dum which used for debugging tool.
	router.Use(gin.Recovery(), middlewares.Logger(), dump.Dump())

	//let we load static files
	router.Static("/css", "./templates/css")
	router.LoadHTMLGlob("./templates/*.html")

	//apiRoute group is used to see/access api via endpoints.
	apiRoute := router.Group("/api", middlewares.BasicAuth())
	{
		apiRoute.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"Message": "hello",
			})
		})

		apiRoute.GET("/videos", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, videocontroller.FindAll())
		})

		apiRoute.POST("/videos", func(ctx *gin.Context) {
			err := videocontroller.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Video input is valid!"})
			}
		})
	}

	//viewRoute Group will use to render static files-no need of authentication
	viewRoute := router.Group("/view")
	{
		viewRoute.GET("/videos", videocontroller.ShowAll)
	}

	//we can get port # from env variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// configOutput create a custom logger file to see debugging outputs.
func configOutput() {
	writer, err := os.Create("simpleAPI.log")
	if err != nil {
		fmt.Printf("unable to create log file")
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, writer)
}
