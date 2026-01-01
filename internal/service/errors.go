package service

import "errors"

var (
	ErrEmailAlreadyInUse   = errors.New("email is already in use")
	ErrDuplicateUsername   = errors.New("username is already in use")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInternalServer      = errors.New("internal server error")
	ErrTokenInvalid        = errors.New("token is invalid")
	ErrTokenExpired        = errors.New("token has expired")
	ErrTokenSignature      = errors.New("token signature is invalid")
	ErrInvalidClaims       = errors.New("token claims are invalid")
	ErrInvalidRefreshToken = errors.New("refresh token is invalid")
)
