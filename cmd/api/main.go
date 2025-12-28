package main

import (
	"log"

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

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		logger.Info("Health check endpoint hit")
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	router.Run(":8080")

}
