package dto

type ErrorResponse struct {
	Error   string `json:"error" example:"Validation failed"`
	Details string `json:"details,omitempty" example:"'Password' failed on the 'min' tag"`
}
