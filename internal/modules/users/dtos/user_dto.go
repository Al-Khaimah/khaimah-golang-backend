package users

import (
	categoryDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/dtos"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	podcastDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/enums"
)

// SignupRequestDTO defines the body for creating a new user.
// swagger:model SignupRequestDTO
type SignupRequestDTO struct {
	FirstName  string   `json:"first_name" validate:"required" example:"John"`
	LastName   string   `json:"last_name" validate:"omitempty" example:"Doe"`
	Email      string   `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Categories []string `json:"categories" validate:"required"  example:"[\"f6a75a87-d695-4ee9-a095-5a79edce4eb8\",\"1db71f60-f276-431d-934c-fd84f8014566\"]"`
	Password   string   `json:"password" validate:"required,passwordvalidator" message:"Password must be at least 6 characters and contain both letters and numbers" example:"Pa55word"`
}

// SignupResponseDTO is returned after a successful signup.
// swagger:model SignupResponseDTO
type SignupResponseDTO struct {
	ID         string                `json:"id" example:"abcd1234"`
	FirstName  string                `json:"first_name" example:"John"`
	LastName   string                `json:"last_name" example:"Doe"`
	Email      string                `json:"email" example:"john.doe@example.com"`
	Categories []categories.Category `json:"categories"`
	Token      string                `json:"token" example:"eyJhbGciOiJIUzI1Ni..."`
	ExpiresAt  string                `json:"expires_at" example:"2025-06-05T15:04:05Z"`
}

// LoginRequestDTO defines the body for logging in.
// swagger:model LoginRequestDTO
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"Pa55word"`
}

// LoginResponseDTO is returned after a successful login.
// swagger:model LoginResponseDTO
type LoginResponseDTO struct {
	ID         string                `json:"id" example:"abcd1234"`
	FirstName  string                `json:"first_name" example:"John"`
	LastName   string                `json:"last_name" example:"Doe"`
	Email      string                `json:"email" example:"john.doe@example.com"`
	Categories []categories.Category `json:"categories"`
	Token      string                `json:"token" example:"eyJhbGciOiJIUzI1Ni..."`
	ExpiresAt  string                `json:"expires_at" example:"2025-06-05T15:04:05Z"`
}

// UserProfileDTO represents a user's profile data.
// swagger:model UserProfileDTO
type UserProfileDTO struct {
	ID         string                `json:"id" example:"abcd1234"`
	FirstName  string                `json:"first_name" example:"John"`
	LastName   string                `json:"last_name" example:"Doe"`
	Email      string                `json:"email" example:"john.doe@example.com"`
	Categories []categories.Category `json:"categories" validate:"required"`
}

// UpdateProfileDTO defines the body for updating a user's profile.
// swagger:model UpdateProfileDTO
type UpdateProfileDTO struct {
	FirstName string `json:"first_name" validate:"omitempty" example:"Johnny"`
	LastName  string `json:"last_name" validate:"omitempty" example:"Doey"`
}

// ChangePasswordDTO defines the body for changing password.
// swagger:model ChangePasswordDTO
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required" example:"OldPa55"`
	NewPassword string `json:"new_password" validate:"required,min=6" example:"NewPa55word"`
}

// UpdatePreferencesDTO defines the body for updating preferences.
// swagger:model UpdatePreferencesDTO
type UpdatePreferencesDTO struct {
	Categories []string `json:"categories" validate:"required" example:"[\"sports\",\"music\"]"`
}

// UpdatePreferencesResponseDTO is returned after preferences are updated.
// swagger:model UpdatePreferencesResponseDTO
type UpdatePreferencesResponseDTO struct {
	Categories categoryDTO.Category `json:"categories"`
}

// GetUserBookmarksResponseDTO wraps a list of bookmarked podcasts.
// swagger:model GetUserBookmarksResponseDTO
type GetUserBookmarksResponseDTO struct {
	Podcasts []podcastDTO.PodcastDto `json:"podcasts"`
}

// CreateSSOUserRequestDTO defines the body for creating an SSO user.
// swagger:model CreateSSOUserRequestDTO
type CreateSSOUserRequestDTO struct {
	FirstName  string   `json:"first_name" validate:"required" example:"Alice"`
	Categories []string `json:"categories" validate:"required" example:"[\"tech\",\"news\"]"`
}

// CreateSSOUserResponseDTO is returned after creating an SSO user.
// swagger:model CreateSSOUserResponseDTO
type CreateSSOUserResponseDTO struct {
	ID         string                `json:"id" example:"xyz789"`
	FirstName  string                `json:"first_name" example:"Alice"`
	Email      string                `json:"email" example:"alice@example.com"`
	Categories []categories.Category `json:"categories"`
}

type UserBaseDTO struct {
	ID         string                `json:"id"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	UserType   users.UserType        `json:"user_type"`
	Email      string                `json:"email"`
	Mobile     string                `json:"mobile"`
	Categories []categories.Category `json:"categories"`
}
