package service

import "errors"

var (
	ErrEmailAlreadyInUse  = errors.New("email is already in use")
	ErrDuplicateUsername  = errors.New("username is already in use")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInternalServer     = errors.New("internal server error")
)
