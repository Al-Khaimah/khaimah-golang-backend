package users

import (
	"fmt"

	"gorm.io/gorm"

	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) CreateUserAuth(userAuth *models.IamAuth) error {
	result := r.DB.Create(userAuth)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}
