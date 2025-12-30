package handler

import (
	"errors"
	"net/http"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/dto"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.Logger
	authService service.AuthInterface
	validator   *validator.Validate
}

func NewAuthHandler(logger *zap.Logger, authService service.AuthInterface, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
		validator:   validator,
	}
}

// Register handles the creation of a new user account.
// @Summary      Register a new user
// @Description  Creates a new user account with a username, email, and password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      dto.SignupRequest  true  "User Registration Details"
// @Success      201   {object}  dto.SignupResponse     "User created successfully"
// @Failure      400   {object}  dto.ErrorResponse    "Invalid input"
// @Failure      409   {object}  dto.ErrorResponse    "Conflict"
// @Failure      500   {object}  dto.ErrorResponse    "Internal Server Error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request format"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("User registration validation failed", zap.Error(err))
		// You can add more sophisticated error parsing here to return specific field errors.
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	user, err := h.authService.RegisterUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyInUse):
			// Use the new struct
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "A user with this email already exists."})
		case errors.Is(err, service.ErrDuplicateUsername):
			// Use the new struct
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "This username is already taken."})
		default:
			h.logger.Error("An unhandled error occurred during registration", zap.Error(err))
			// Use the new struct
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "An internal error occurred."})
		}
		return
	}

	resp := dto.SignupResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles user authentication and token issuance.
// @Summary      User login
// @Description  Authenticates a user and issues access and refresh tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginRequest  true  "User Login Credentials"
// @Success      200          {object}  dto.LoginResponse     "Login successful"
// @Failure      400          {object}  dto.ErrorResponse    "Invalid input"
// @Failure      401          {object}  dto.ErrorResponse    "Unauthorized"
// @Failure      500          {object}  dto.ErrorResponse    "Internal Server Error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request format"})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "User not found"})

		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid email or password"})
		default:
			h.logger.Error("An unhandled error occurred during login", zap.Error(err))
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "An internal error occurred."})
		}
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		Domain:   "localhost",
		MaxAge:   3600 * 24 * 7, // 1 week in seconds
		Secure:   false,         // Set to true in production with HTTPS
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)

	resp := dto.LoginResponse{
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, resp)

}
