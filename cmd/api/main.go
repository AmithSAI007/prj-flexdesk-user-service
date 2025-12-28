package main

import (
	"log"

	_ "github.com/AmithSAI007/prj-flexdesk-user-service/docs"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/api"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/config"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/handler"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

// @title        FlexDesk User Service API
// @version      1.0
// @description  This is the API for managing users...
// @BasePath  /api/v1
func main() {
	logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.PrometheusMetrics())

	userHandler := handler.NewUserHandler(logger)
	handlers := &api.HandlerRegistry{
		UserHandler: userHandler,
	}
	api.SetupRoutes(router, handlers)

	logger.Info("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		logger.Error(err.Error())
	}
}
