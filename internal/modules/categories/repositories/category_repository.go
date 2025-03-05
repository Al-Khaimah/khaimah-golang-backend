package categories

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

func (r *CategoryRepository) FindAllCategories() ([]models.Category, error) {
	var categories []models.Category
	result := r.DB.Find(&categories)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", result.Error)
	}
	return categories, nil
}

func (r *CategoryRepository) FindCategoryByID(categoryID uuid.UUID) (*models.Category, error) {
	var category models.Category
	result := r.DB.Where("id = ?", categoryID).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find category: %w", result.Error)
	}
	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	result := r.DB.Create(category)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create category: %w", result.Error)
	}
	return category, nil
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	result := r.DB.Save(category)
	if result.Error != nil {
		return fmt.Errorf("failed to update category: %w", result.Error)
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(categoryID uuid.UUID) error {
	result := r.DB.Delete(&models.Category{}, categoryID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}
	return nil
}
