package main

import (
	"fmt"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/glue"
	middlewares2 "gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/middleware"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/rest"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/repository/gorm"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/hexagonalarch/controller"
	"gitlab.com/Yinebeb-01/hexagonalarch/services"

	dump "github.com/tpkeeper/gin-dump"
)

var (
	videoRepository = gorm.NewVideoRepository()
	loginService    = services.NewLoginService()
	jwtService      = services.NewJWTService()

	videoService    = service.New(videoRepository)
	videoHandler    = rest.Init(videoService)
	loginController = controller.NewLoginController(loginService, jwtService)
)

func main() {
	configOutput()

	router := gin.New()
	// middleware dump is an alise of gin-dum which used for debugging tool.
	router.Use(gin.Recovery(), middlewares2.Logger(), dump.Dump())
	router.Static("/css", "./internal/adapter/templates/css")
	router.LoadHTMLGlob("./internal/adapter/templates/*.html")

	// Login Endpoint: Authentication
	router.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, nil)
		}
	})

	//viewRoute Group will use to render static files
	viewRoute := router.Group("/view")
	{
		viewRoute.GET("/videos", videoHandler.ShowAll)
	}

	v1 := router.Group("/v1")
	glue.InitVideoRoute(v1, videoHandler)
	log.Println("router initialized")

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
