package users

import categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"

type OAuthRequestDTO struct {
	Provider string `json:"provider"`
}

type UserDetailsDTO struct {
	ID         string                `json:"id" example:"xyz789"`
	FirstName  string                `json:"first_name" example:"Alice"`
	Email      string                `json:"email" example:"alice@example.com"`
	Categories []categories.Category `json:"categories"`
}
type SSOLoginResponse struct {
	Token      string          `json:"token"`
	UserExists bool            `json:"user_exists"`
	User       *UserDetailsDTO `json:"user"`
}
