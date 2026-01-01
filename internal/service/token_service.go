package service

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/db"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenInterface interface {
	NewTokenPair(user *model.User, now time.Time) (accessToken string, refreshToken string, err error)
	ValidateAccessToken(tokenStr string) (*TokenClaims, error)
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type TokenService struct {
	logger          *zap.Logger
	issuer          string
	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	accessDuration  time.Duration
	refreshDuration time.Duration
	store           db.Store
}

func NewTokenService(issuer string, logger *zap.Logger, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, store db.Store) TokenInterface {

	return &TokenService{
		issuer:          issuer,
		accessDuration:  time.Minute * 15,
		refreshDuration: time.Hour * 24 * 7,
		logger:          logger,
		privateKey:      privateKey,
		publicKey:       publicKey,
		store:           store,
	}
}

var _ TokenInterface = (*TokenService)(nil)

func (s *TokenService) NewTokenPair(user *model.User, now time.Time) (string, string, error) {

	accessToken, err := s.generateToken(user.ID.String(), user.Email, TokenTypeAccess, now)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user.ID.String(), user.Email, TokenTypeRefresh, now)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) ValidateAccessToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrTokenExpired):
			s.logger.Error("Token has expired", zap.Error(err))
			return nil, ErrTokenExpired
		case errors.Is(err, ErrTokenSignature):
			s.logger.Error("Token signature is invalid", zap.Error(err))
			return nil, ErrTokenSignature
		default:
			s.logger.Error("Failed to parse token", zap.Error(err))
			return nil, ErrTokenInvalid
		}

	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.Type != TokenTypeAccess {
			s.logger.Error("Invalid token type", zap.String("type", claims.Type))
			return nil, ErrInvalidClaims
		}
		return claims, nil
	} else {
		s.logger.Error("Invalid token claims")
		return nil, ErrInvalidClaims
	}
}

func (s *TokenService) generateToken(userID, email, tokenType string, now time.Time) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		s.logger.Error("Failed to sign token", zap.Error(err))
		return "", fmt.Errorf("failed to sign token: %w ", err)
	}

	s.logger.Info("Token generated successfully for user", zap.String("user_id", userID), zap.String("type", tokenType))
	return signedToken, nil

}
