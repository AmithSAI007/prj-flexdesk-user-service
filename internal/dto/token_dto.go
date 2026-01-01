package dto

type TokenDto struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	TokenHash string `json:"token_hash"`
	IsActive  bool   `json:"is_active"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}
