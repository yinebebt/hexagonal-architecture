package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yinebebt/hexagonal-architecture/docs"
	gqlhandler "github.com/yinebebt/hexagonal-architecture/internal/adapter/graphql"
	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository"
	"github.com/yinebebt/hexagonal-architecture/internal/adapter/rest"
	"github.com/yinebebt/hexagonal-architecture/internal/core/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	dbType = flag.String("dbtype", "sqlite", "Type of the database (e.g. sqlite,postgres)")
	dsn    = flag.String("dsn", "test.db", "Data source name for the database")
)

// @title			hexagonal-architecture
// @version		0.1.0
// @contact.name	Yinebe T.
// @contact.url	www.linkedin.com/in/yinebeb-tariku
// @contact.email	yintar5@gmail.com
// @host			localhost
// @BasePath		/v1
func main() {
	flag.Parse()

	videoRepository, err := repository.NewVideoRepository(*dbType, *dsn)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer func() {
		if err := videoRepository.Close(); err != nil {
			log.Printf("Error closing repository: %v", err)
		}
	}()

	videoService := service.New(videoRepository)
	videoHandler := rest.NewVideoHandler(videoService)

	logWriter, err := configOutput()
	if err != nil {
		log.Printf("Warning: failed to setup log file: %v", err)
	} else if logWriter != nil {
		defer func() {
			if err := logWriter.Close(); err != nil {
				log.Printf("Error closing log file: %v", err)
			}
		}()
	}

	router := gin.New()
	router.Use(gin.Recovery(), rest.Logger())
	router.Static("/css", "./internal/adapter/templates/css")
	router.LoadHTMLGlob("./internal/adapter/templates/*.html")

	v1 := router.Group("/v1")
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rest.RegisterVideoRoutes(v1, videoHandler)

	// GraphQL adapter — serves playground at /v1/graphql
	graphqlHandler, err := gqlhandler.NewHandler(videoService)
	if err != nil {
		log.Fatalf("Failed to initialize GraphQL handler: %v", err)
	}
	v1.Any("/graphql", gin.WrapH(graphqlHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// configOutput creates a custom logger file to see debugging outputs.
// Returns the file writer so it can be closed properly.
func configOutput() (io.Closer, error) {
	writer, err := os.Create("app.log")
	if err != nil {
		return nil, fmt.Errorf("unable to create log file: %w", err)
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, writer)
	return writer, nil
}
