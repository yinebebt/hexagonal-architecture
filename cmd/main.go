package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/Yinebeb-01/hexagonalarch/docs"
	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/glue/route"
	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/middleware"
	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/handler/rest"
	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/repository"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	dump "github.com/tpkeeper/gin-dump"
)

var (
	dbType = flag.String("dbtype", "sqlite", "Type of the database (sqliteor postgres)")
	dsn    = flag.String("dsn", "test.db", "Data source name for the database")
)

// @title hexagonal-architecture
// @version         0.1.0
// @contact.name   Yinebe T.
// @contact.url    www.linkedin.com/in/yinebeb-tariku
// @contact.email  yintar5@gmail.com
// @host localhost
// @BasePath  /v1
func main() {
	flag.Parse()

	videoRepository := repository.NewVideoRepository(*dbType, *dsn)
	loginService := service.NewLoginService()
	jwtService := service.NewJWTService()

	videoService := service.New(videoRepository)
	videoHandler := rest.InitVideo(videoService)
	loginHandler := rest.InitLogin(loginService, jwtService)

	configOutput()

	router := gin.New()
	// dump is an alise of gin-dum, used for debugging.
	router.Use(gin.Recovery(), middleware.Logger(), dump.Dump())
	router.Static("/css", "./internal/adapter/templates/css")
	router.LoadHTMLGlob("./internal/adapter/templates/*.html")

	v1 := router.Group("/v1")
	docs.SwaggerInfo.Host = "server.host"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"swagger.schemes"}
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	writer, err := os.Create("app.log")
	if err != nil {
		log.Println("unable to create log file")
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, writer)
}
