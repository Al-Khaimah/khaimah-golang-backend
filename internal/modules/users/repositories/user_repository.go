package users

import (
	"errors"
	"fmt"

	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindOneByEmail(email string) (*models.User, error) {
	var user models.User

	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	result := r.DB.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %w", result.Error)
	}

	return user, nil
}
