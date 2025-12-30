package handler

import (
	"errors"
	"net/http"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/constants"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/dto"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.Logger
	userService service.UserInterface
	validator   *validator.Validate
}

func NewUserHandler(logger *zap.Logger, userService service.UserInterface, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
		validator:   validator,
	}
}

// GetUserProfile handles fetching the profile of the authenticated user.
// @Summary      Get User Profile
// @Description  Retrieves the profile information of the authenticated user.
// @Tags         User
// @Produce      json
// @Success      200  {object}  dto.UserProfileResponse  "User profile retrieved successfully"
// @Failure      401  {object}  dto.ErrorResponse       "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse       "User not found"
// @Failure      500  {object}  dto.ErrorResponse       "Internal Server Error"
// @Security     BearerAuth
// @Router       /user/profile [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Implementation for getting user profile
	userIDVal, exists := c.Request.Context().Value(constants.UserIDKey).(string)
	if !exists {
		h.logger.Error("UserID not found in context, middleware might be missing")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "An internal error occurred"})
		return
	}

	userId, err := uuid.Parse(userIDVal)
	if err != nil {
		h.logger.Error("Invalid UserID format", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			h.logger.Warn("User not found", zap.String("userID", userIDVal))
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		default:
			h.logger.Error("Failed to get user profile", zap.Error(err))
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "An internal error occurred"})
		}
	}

	resp := dto.UserProfileResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}
