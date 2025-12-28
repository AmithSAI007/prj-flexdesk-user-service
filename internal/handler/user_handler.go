package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger *zap.Logger
}

func NewUserHandler(logger *zap.Logger) *UserHandler {
	return &UserHandler{
		logger: logger,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	h.logger.Info("CreateUser handler called")
	// Implementation for creating a user goes here
}
