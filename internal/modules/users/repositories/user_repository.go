package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"

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

func (r *UserRepository) FindOneByID(userID uuid.UUID) (*models.User, error) {
	var user models.User

	result := r.DB.Where("id = ?", userID).First(&user)
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

func (r *UserRepository) UpdateUser(user *models.User) error {
	result := r.DB.Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user profile: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := r.DB.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", result.Error)
	}
	return users, nil
}

func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	result := r.DB.Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	return nil
}
