package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/model"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/security"

	"go.uber.org/zap"
)

type AuthInterface interface {
	// Define user-related methods here, e.g., CreateUser, GetUser, UpdateUser, DeleteUser, etc.
	RegisterUser(ctx context.Context, username, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	persistRefreshToken(ctx context.Context, userID string, refreshToken string, expiration time.Time) error
}

type AuthService struct {
	logger         *zap.Logger
	userInterface  UserInterface
	tokenInterface TokenInterface
}

func NewAuthService(logger *zap.Logger, userInterface UserInterface, tokenInterface TokenInterface) AuthInterface {
	return &AuthService{
		logger:         logger,
		userInterface:  userInterface,
		tokenInterface: tokenInterface,
	}
}

var _ AuthInterface = (*AuthService)(nil)

func (s *AuthService) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	// Implement user registration logic here
	// For example, hash the password, validate input, and store the user in the database

	_, err := s.userInterface.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyInUse
	}

	hashedPassword, err := security.HashPassword(password) // Implement password hashing

	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, ErrInternalServer
	}

	user, err := s.userInterface.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {

	user, err := s.userInterface.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	err = security.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		s.logger.Warn("Invalid password attempt", zap.String("email", email))
		return "", "", ErrInvalidCredentials
	}

	now := time.Now()

	accessToken, refreshToken, err := s.tokenInterface.NewTokenPair(&model.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, now)
	if err != nil {
		s.logger.Error("Failed to generate token pair", zap.Error(err))
		return "", "", ErrInternalServer
	}

	expiration := now.Add(7 * 24 * time.Hour) // Assuming refresh token validity is 7 days

	err = s.persistRefreshToken(ctx, user.ID.String(), refreshToken, expiration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) persistRefreshToken(ctx context.Context, userID string, refreshToken string, expiration time.Time) error {
	// Implement logic to persist refresh token, e.g., store in database or cache
	hash := sha256.New()
	hash.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash.Sum(nil))

	tokenErr := s.tokenInterface.SaveRefreshToken(ctx, userID, tokenHash, expiration)
	if tokenErr != nil {
		s.logger.Error("Failed to persist refresh token", zap.Error(tokenErr))
		return ErrInternalServer
	}

	return nil
}
