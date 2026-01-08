package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/db"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/model"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/security"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	generated "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"

	"go.uber.org/zap"
)

type AuthInterface interface {
	// Define user-related methods here, e.g., CreateUser, GetUser, UpdateUser, DeleteUser, etc.
	RegisterUser(ctx context.Context, username, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthService struct {
	logger         *zap.Logger
	store          db.Store
	userInterface  UserInterface
	tokenInterface TokenInterface
}

func NewAuthService(logger *zap.Logger, userInterface UserInterface, tokenInterface TokenInterface, store db.Store) AuthInterface {
	return &AuthService{
		logger:         logger,
		userInterface:  userInterface,
		tokenInterface: tokenInterface,
		store:          store,
	}
}

var _ AuthInterface = (*AuthService)(nil)

func (s *AuthService) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	// Implement user registration logic here
	// For example, hash the password, validate input, and store the user in the database
	var user *model.User

	err := s.store.ExecTx(ctx, func(q generated.Querier) error {

		_, err := s.userInterface.GetUserByEmail(ctx, q, email)
		if err == nil {
			return ErrEmailAlreadyInUse
		}

		hashedPassword, err := security.HashPassword(password) // Implement password hashing

		if err != nil {
			s.logger.Error("Failed to hash password", zap.Error(err))
			return ErrInternalServer
		}

		createdUser, err := s.userInterface.CreateUser(ctx, q, username, email, hashedPassword)
		if err != nil {
			return err
		}

		user = createdUser

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {

	var accessToken, refreshToken string

	err := s.store.ExecTx(ctx, func(q generated.Querier) error {

		user, err := s.userInterface.GetUserByEmail(ctx, q, email)
		if err != nil {
			return ErrUserNotFound
		}

		err = security.CheckPasswordHash(password, user.PasswordHash)
		if err != nil {
			s.logger.Warn("Invalid password attempt", zap.String("email", email))
			return ErrInvalidCredentials
		}

		gAccessToken, gRefreshToken, err := s.issueNewTokenPair(ctx, q, user)
		if err != nil {
			return err
		}

		accessToken = gAccessToken
		refreshToken = gRefreshToken

		return nil
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, rfToken string) (string, string, error) {

	hash := sha256.New()
	hash.Write([]byte(rfToken))
	tokenHash := hex.EncodeToString(hash.Sum(nil))
	var accessToken, refreshToken string

	err := s.store.ExecTx(ctx, func(q generated.Querier) error {

		token, err := q.GetRefreshToken(context.Background(), tokenHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				s.logger.Error("Refresh token not found", zap.String("token_hash", tokenHash))
				return ErrInvalidRefreshToken
			}
			s.logger.Error("Failed to retrieve refresh token", zap.Error(err))
			return ErrInternalServer
		}

		err = q.InvalidateRefreshToken(ctx, token.ID)
		if err != nil {
			s.logger.Error("Failed to invalidate refresh token", zap.Error(err), zap.String("tokenID", token.ID.String()))
			return ErrInternalServer
		}

		if !token.IsActive || token.ExpiresAt.Time.Before(time.Now()) {
			s.logger.Warn("Refresh token is inactive or expired", zap.String("tokenHash", tokenHash))
			return ErrInvalidRefreshToken
		}

		user, err := s.userInterface.GetUserByID(ctx, q, token.UserID.Bytes)
		if err != nil {
			s.logger.Warn("User associated with refresh token not found", zap.String("userID", token.UserID.String()))
			return ErrInvalidRefreshToken
		}

		gAccessToken, gAefreshToken, err := s.issueNewTokenPair(ctx, q, user)
		if err != nil {
			return err
		}
		accessToken = gAccessToken
		refreshToken = gAefreshToken

		return nil

	})

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func (s *AuthService) hashRefreshToken(token string) string {
	// Persist the NEW refresh token
	hash := sha256.New()
	hash.Write([]byte(token))
	tokenHash := hex.EncodeToString(hash.Sum(nil))
	return tokenHash

}

func (s *AuthService) issueNewTokenPair(ctx context.Context, q generated.Querier, user *model.User) (string, string, error) {
	now := time.Now()
	accessToken, refreshToken, err := s.tokenInterface.NewTokenPair(user, now)
	if err != nil {
		s.logger.Error("Failed to generate token pair", zap.Error(err), zap.String("userID", user.ID.String()))
		return "", "", ErrInternalServer
	}

	expiration := now.Add(7 * 24 * time.Hour) // 7 day validity for refresh token

	tokenHash := s.hashRefreshToken(refreshToken)

	_, err = q.CreateRefreshToken(ctx, generated.CreateRefreshTokenParams{
		UserID:    pgtype.UUID{Bytes: user.ID, Valid: true},
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiration, Valid: true},
	})
	if err != nil {
		s.logger.Error("Failed to persist new refresh token", zap.Error(err), zap.String("userID", user.ID.String()))
		return "", "", ErrInternalServer
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {

	tokenHash := s.hashRefreshToken(refreshToken)

	err := s.store.InvalidateRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		s.logger.Error("Failed to invalidate refresh token during logout", zap.Error(err), zap.String("tokenHash", tokenHash))
		return ErrInternalServer
	}

	return nil
}
