package categories

import (
	base "github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categoryDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/dtos"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	categoryRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"
	"github.com/google/uuid"
)

type CategoryService struct {
	CategoryRepo *categoryRepository.CategoryRepository
}

func NewCategoryService(categoryRepo *categoryRepository.CategoryRepository) *CategoryService {
	return &CategoryService{CategoryRepo: categoryRepo}
}

func (s *CategoryService) GetAllCategories() base.Response {
	categories, err := s.CategoryRepo.FindAllCategories()
	if err != nil {
		return base.SetErrorMessage("Failed to fetch categories", err)
	}
	if len(categories) == 0 {
		return base.SetData([]categoryDTO.Category{}, "No categories found")
	}

	var categoryResponse []interface{}
	for _, category := range categories {
		categoryResponse = append(categoryResponse, categoryDTO.Category{
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return base.SetData(categoryResponse)
}

func (s *CategoryService) CreateCategory(categoryData categoryDTO.Category) base.Response {
	newCategory := &models.Category{
		Name:            categoryData.Name,
		Description:     categoryData.Description,
		IsNewsIntensive: categoryData.IsNewsIntensive,
	}

	createdCategory, err := s.CategoryRepo.CreateCategory(newCategory)
	if err != nil {
		return base.SetErrorMessage("Failed to create category", err)
	}

	if createdCategory == nil {
		return base.SetErrorMessage("Failed to create category", "Category creation returned nil")
	}

	categoryResponse := categoryDTO.Category{
		ID:              createdCategory.ID.String(),
		Name:            createdCategory.Name,
		Description:     createdCategory.Description,
		IsNewsIntensive: createdCategory.IsNewsIntensive,
	}

	return base.SetData(categoryResponse, "Category created successfully")
}

func (s *CategoryService) UpdateCategory(categoryID string, updateData categoryDTO.Category) base.Response {
	uid, err := uuid.Parse(categoryID)
	if err != nil {
		return base.SetErrorMessage("Invalid Category ID", err)
	}

	category, err := s.CategoryRepo.FindCategoryByID(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to fetch category", err)
	}
	if category == nil {
		return base.SetErrorMessage("Category not found", "No category exists with this ID")
	}

	if updateData.Name != "" {
		category.Name = updateData.Name
	}
	if updateData.Description != "" {
		category.Description = updateData.Description
	}
	category.IsNewsIntensive = updateData.IsNewsIntensive

	err = s.CategoryRepo.UpdateCategory(category)
	if err != nil {
		return base.SetErrorMessage("Failed to update category", err)
	}

	categoryResponse := categoryDTO.Category{
		ID:              category.ID.String(),
		Name:            category.Name,
		Description:     category.Description,
		IsNewsIntensive: category.IsNewsIntensive,
	}

	return base.SetData(categoryResponse, "Category updated successfully")
}

func (s *CategoryService) DeleteCategory(categoryID string) base.Response {
	uid, err := uuid.Parse(categoryID)
	if err != nil {
		return base.SetErrorMessage("Invalid Category ID", err)
	}

	category, err := s.CategoryRepo.FindCategoryByID(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to fetch category", err)
	}
	if category == nil {
		return base.SetErrorMessage("Category not found", "No category exists with this ID")
	}

	err = s.CategoryRepo.DeleteCategory(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to delete category", err)
	}

	return base.SetSuccessMessage("Category deleted successfully")
}
