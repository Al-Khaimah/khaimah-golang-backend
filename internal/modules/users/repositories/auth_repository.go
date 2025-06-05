package users

import (
	"fmt"

	"github.com/google/uuid"

	"gorm.io/gorm"

	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) CreateUserAuth(userAuth *models.IamAuth) error {
	result := r.DB.Create(userAuth)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (r *AuthRepository) FindAuthByUserID(userID uuid.UUID) (*models.IamAuth, error) {
	var auth models.IamAuth
	result := r.DB.Where("user_id = ?", userID).First(&auth)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find authentication record: %w", result.Error)
	}
	return &auth, nil
}

func (r *AuthRepository) FindAuthByOTP(otp int) (*models.IamAuth, error) {
	var auth models.IamAuth
	result := r.DB.Where("otp = ?", otp).First(&auth)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find authentication record: %w", result.Error)
	}
	return &auth, nil
}

func (r *AuthRepository) UpdateAuth(auth *models.IamAuth) error {
	result := r.DB.Save(auth)
	if result.Error != nil {
		return fmt.Errorf("failed to update authentication record: %w", result.Error)
	}
	return nil
}
