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

// @BasePath /api/v1
// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      map[string]interface{}  true  "User Data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.logger.Info("CreateUser handler called")
	// Implementation for creating a user goes here
}
