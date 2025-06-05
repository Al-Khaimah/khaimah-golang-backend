package users

import (
	"errors"
	"fmt"

	podcastModel "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"

	"github.com/google/uuid"

	categoryModel "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
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

	result := r.DB.
		Preload("Categories").
		Preload("Bookmarks").
		Preload("Downloads").
		Where("email ILIKE ?", email).
		First(&user)

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

	result := r.DB.
		Preload("Categories").
		Preload("Bookmarks").
		Preload("Downloads").
		Where("id = ?", userID).First(&user)
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

func (r *UserRepository) FindAllUsers() ([]models.User, error) {
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

func (r *UserRepository) FindUserCategories(userID uuid.UUID) ([]categoryModel.Category, error) {
	var user models.User

	result := r.DB.Preload("Categories").First(&user, "id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", result.Error)
	}

	return user.Categories, nil
}

func (r *UserRepository) FindDownloadedPodcasts(userID uuid.UUID) ([]podcastModel.Podcast, error) {
	var user models.User

	result := r.DB.Where("id = ?", userID).Preload("Downloads").First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return []podcastModel.Podcast{}, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return user.Downloads, nil
}

func (r *UserRepository) UpdateUserPreferences(user *models.User, categories []categoryModel.Category) error {
	err := r.DB.Model(user).Association("Categories").Replace(categories)
	if err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}
	return nil
}

func (r *UserRepository) FindOrCreateByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return &user, nil
	}
	user = models.User{Email: email}
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindOneByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User

	result := r.DB.
		Preload("Categories").
		Preload("Bookmarks").
		Preload("Downloads").
		Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return &user, nil
}
