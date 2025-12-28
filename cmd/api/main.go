package main

import (
	"log"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/api"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/handler"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "app.log"}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))

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
