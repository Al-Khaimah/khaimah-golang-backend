package users

import (
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
)

// SendOTPRequestDTO defines the body for sending an OTP.
// swagger:model SendOTPRequestDTO
type SendOTPRequestDTO struct {
	Mobile     string   `json:"mobile" validate:"omitempty" example:"+9665XXXXXXX"`
	Email      string   `json:"email" validate:"omitempty,email" example:"user@example.com"`
	FirstName  string   `json:"first_name" validate:"omitempty" example:"Ziyad"`
	Categories []string `json:"categories" validate:"omitempty" example:"[\"9d671bac-17b0-42cf-b68a-aa908f30b134\"]"`
}

// SendOTPResponseDTO is returned after successfully sending an OTP.
// swagger:model SendOTPResponseDTO
type SendOTPResponseDTO struct {
	Message string `json:"message" example:"OTP sent successfully"`
}

// VerifyOTPRequestDTO defines the body for verifying an OTP.
// swagger:model VerifyOTPRequestDTO
type VerifyOTPRequestDTO struct {
	Mobile string `json:"mobile" validate:"omitempty" example:"+9665XXXXXXX"`
	Email  string `json:"email" validate:"omitempty,email" example:"user@example.com"`
	OTP    string `json:"otp" validate:"required" example:"1234"`
}

// VerifyOTPResponseDTO is returned after successfully verifying an OTP.
// swagger:model VerifyOTPResponseDTO
type VerifyOTPResponseDTO struct {
	ID         string                `json:"id" example:"abcd1234"`
	FirstName  string                `json:"first_name" example:"Ziyad"`
	Email      string                `json:"email,omitempty" example:"user@example.com"`
	Mobile     string                `json:"mobile,omitempty" example:"+9665XXXXXXX"`
	Categories []categories.Category `json:"categories"`
	Token      string                `json:"token" example:"eyJhbGciOiJIUzI1Ni..."`
	ExpiresAt  string                `json:"expires_at" example:"2025-06-05T15:04:05Z"`
}
