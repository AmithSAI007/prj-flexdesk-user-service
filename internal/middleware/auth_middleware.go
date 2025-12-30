package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/constants"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	logger       *zap.Logger
	tokenService service.TokenInterface
}

func NewAuthMiddleware(logger *zap.Logger, tokenService service.TokenInterface) *AuthMiddleware {
	return &AuthMiddleware{
		logger:       logger,
		tokenService: tokenService,
	}
}

func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement authentication logic here
		// For example, validate JWT tokens, check user roles, etc.}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			a.logger.Warn("Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.logger.Warn("Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}

		tokenStr := parts[1]

		claims, err := a.tokenService.ValidateAccessToken(tokenStr)
		if err != nil {
			a.logger.Warn("Invalid token", zap.Error(err))
			switch {
			case errors.Is(err, service.ErrTokenExpired):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token_expired"})
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			}

			return
		}

		ctx := context.WithValue(c.Request.Context(), constants.UserIDKey, claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

	}
}
