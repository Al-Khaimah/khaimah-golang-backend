package users

import categoryDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/dtos"

type SignupRequestDTO struct {
	FirstName  string   `json:"first_name" validate:"required"`
	LastName   string   `json:"last_name" validate:"omitempty"`
	Email      string   `json:"email" validate:"required,email"`
	Categories []string `json:"categories" validate:"omitempty"`
	Password   string   `json:"password" validate:"required"`
}

type SignupResponseDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type UserProfileDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateProfileDTO struct {
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdatePreferencesDTO struct {
	Categories []string `json:"categories" validate:"required"`
}

type UpdatePreferencesResponseDTO struct {
	Categories categoryDTO.Category `json:"categories" validate:"required"`
}
