package api

import (
	"github.com/AmithSAI007/prj-flexdesk-user-service/docs"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))

	router.GET("/docs/doc.json", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.WriteHeader(200)
		ctx.Writer.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
