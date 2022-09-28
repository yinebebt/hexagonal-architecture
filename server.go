package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/simpleAPI/controller"
	"gitlab.com/Yinebeb-01/simpleAPI/middlewares"
	"gitlab.com/Yinebeb-01/simpleAPI/repository"
	"gitlab.com/Yinebeb-01/simpleAPI/services"

	dump "github.com/tpkeeper/gin-dump"
)

var (
	videorepository repository.VideoReposiory = repository.NewVideoRepository()
	lognservice     services.LoginService     = services.NewLoginService()
	jwtservice      services.JWTService       = services.NewJWTService()

	videoservice    services.VideoService      = services.New(videorepository)
	videocontroller controller.VideoController = controller.New(videoservice)
	logincontroller controller.LoginController = controller.NewLoginController(lognservice, jwtservice)
)

func main() {
	defer videorepository.Close()
	configOutput()

	router := gin.New()
	//custome middlewares used, dump is an alise of gin-dum which used for debugging tool.
	router.Use(gin.Recovery(), middlewares.Logger(), dump.Dump())

	//let we load static files
	router.Static("/css", "./templates/css")
	router.LoadHTMLGlob("./templates/*.html")

	// Login Endpoint: Authentication + Token creation
	router.POST("/login", func(ctx *gin.Context) {
		token := logincontroller.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	// JWT Authorization Middleware applies to "/api" only. You can see/check the Basc Auth too
	//apiRoute group is used to see/access api via endpoints.
	apiRoute := router.Group("/api", middlewares.AuthorizeJWT()) //middlewares.BasicAuth()
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

		apiRoute.PUT("/videos/:id", func(ctx *gin.Context) {
			err := videocontroller.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Video input is valid!"})
			}
		})

		apiRoute.DELETE("/videos/:id", func(ctx *gin.Context) {
			err := videocontroller.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Video Deleted!"})
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
