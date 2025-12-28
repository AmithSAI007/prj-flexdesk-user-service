package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic", zap.Any("error", r))
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
				})
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					ErrorResponse{
						Status:  http.StatusInternalServerError,
						Message: "Internal Server Error",
					},
				)
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Error("Request error", zap.Error(e.Err))
			}
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
			})
		}

	}
}
