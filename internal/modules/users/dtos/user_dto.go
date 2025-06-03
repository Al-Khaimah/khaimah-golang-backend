package users

import (
	categoryDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/dtos"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	podcastDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
)

type SignupRequestDTO struct {
	FirstName  string   `json:"first_name" validate:"required"`
	LastName   string   `json:"last_name" validate:"omitempty"`
	Email      string   `json:"email" validate:"required,email"`
	Categories []string `json:"categories" validate:"required"`
	Password   string   `json:"password" validate:"required,passwordvalidator" message:"Password must be at least 6 characters and contain both letters and numbers"`
}

type SignupResponseDTO struct {
	ID         string                `json:"id"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Email      string                `json:"email"`
	Categories []categories.Category `json:"categories"`
	Token      string                `json:"token"`
	ExpiresAt  string                `json:"expires_at"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	ID         string                `json:"id"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Email      string                `json:"email"`
	Categories []categories.Category `json:"categories"`
	Token      string                `json:"token"`
	ExpiresAt  string                `json:"expires_at"`
}

type UserProfileDTO struct {
	ID         string                `json:"id"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Email      string                `json:"email"`
	Categories []categories.Category `json:"categories" validate:"required"`
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
	Categories categoryDTO.Category `json:"categories"`
}

type GetUserBookmarksResponseDTO struct {
	Podcasts []podcastDTO.PodcastDto `json:"podcasts"`
}

type CreateSSOUserRequestDTO struct {
	FirstName  string   `json:"first_name" validate:"required"`
	Categories []string `json:"categories" validate:"required"`
}

type CreateSSOUserResponseDTO struct {
	ID         string                `json:"id"`
	FirstName  string                `json:"first_name"`
	Email      string                `json:"email"`
	Categories []categories.Category `json:"categories"`
}
