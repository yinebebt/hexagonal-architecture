package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/hexagonalarch/controller"
	"gitlab.com/Yinebeb-01/hexagonalarch/middlewares"
	"gitlab.com/Yinebeb-01/hexagonalarch/repository"
	"gitlab.com/Yinebeb-01/hexagonalarch/services"

	dump "github.com/tpkeeper/gin-dump"
)

var (
	videoRepository = repository.NewVideoRepository()
	loginService    = services.NewLoginService()
	jwtService      = services.NewJWTService()

	videoService    = services.New(videoRepository)
	videoController = controller.New(videoService)
	loginController = controller.NewLoginController(loginService, jwtService)
)

func main() {
	configOutput()
	router := gin.New()
	// middlewares dump is an alise of gin-dum which used for debugging tool.
	router.Use(gin.Recovery(), middlewares.Logger(), dump.Dump())
	router.Static("/css", "./templates/css")
	router.LoadHTMLGlob("./templates/*.html")

	// Login Endpoint: Authentication + Token creation
	router.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	//apiRoute group used to group 'api/*' endpoints.
	apiRoute := router.Group("/api", middlewares.BasicAuth(), middlewares.AuthorizeJWT())
	{
		apiRoute.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"Message": "hello",
			})
		})

		apiRoute.GET("/videos", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, videoController.FindAll())
		})

		apiRoute.POST("/videos", func(ctx *gin.Context) {
			err := videoController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Video input is valid!"})
			}
		})

		apiRoute.PUT("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Video input is invalid!"})
			}
		})

		apiRoute.DELETE("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Video Deleted!"})
			}
		})
	}

	//viewRoute Group will use to render static files
	viewRoute := router.Group("/view")
	{
		viewRoute.GET("/videos", videoController.ShowAll)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// configOutput create a custom logger file to see debugging outputs.
func configOutput() {
	writer, err := os.Create("hexagonal_arch.log")
	if err != nil {
		fmt.Printf("unable to create log file")
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, writer)
}
