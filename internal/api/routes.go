package api

import (
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/handler"
	"github.com/gin-gonic/gin"
)

type HandlerRegistry struct {
	// Add your handlers here
	UserHandler *handler.UserHandler
}

func SetupRoutes(router *gin.Engine, handlers *HandlerRegistry) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", handlers.UserHandler.CreateUser)
		// Add more user routes as needed
	}
}
