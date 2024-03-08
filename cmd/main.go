package main

import (
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/glue/route"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/middleware"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/rest"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/repository/gorm"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	dump "github.com/tpkeeper/gin-dump"
)

var (
	videoRepository = gorm.NewVideoRepository()
	loginService    = service.NewLoginService()
	jwtService      = service.NewJWTService()

	videoService = service.New(videoRepository)
	videoHandler = rest.InitVideo(videoService)
	loginHandler = rest.InitLogin(loginService, jwtService)
)

func main() {
	configOutput()

	router := gin.New()
	// dump is an alise of gin-dum, used for debugging.
	router.Use(gin.Recovery(), middleware.Logger(), dump.Dump())
	router.Static("/css", "./internal/adapter/templates/css")
	router.LoadHTMLGlob("./internal/adapter/templates/*.html")

	v1 := router.Group("/v1")
	route.InitVideoRoute(v1, videoHandler)
	route.InitLoginRoute(v1, loginHandler)

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
		log.Println("unable to create log file")
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, writer)
}
