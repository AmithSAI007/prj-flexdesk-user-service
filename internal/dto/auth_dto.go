package dto

import (
	"github.com/google/uuid"
	"time"
)

type SignupRequest struct {
	Username string `json:"username" binding:"required" validate:"min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required" validate:"min=6"`
}

type SignupResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
