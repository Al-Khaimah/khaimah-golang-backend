package categories

import (
	base "github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categoryDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/dtos"
	categoryRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"
)

type CategoryService struct {
	CategoryRepository *categoryRepository.CategoryRepository
}

func NewCategoryService(categoryRepository *categoryRepository.CategoryRepository) *CategoryService {
	return &CategoryService{CategoryRepository: categoryRepository}
}

func (s *CategoryService) GetCategories() base.Response {
	categories, err := s.CategoryRepository.FindAllCategories()
	if err != nil {
		return base.SetErrorMessage("Failed to fetch categories", err)
	}

	var categoryDTOs []categoryDTO.Category
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, categoryDTO.Category{
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return base.SetData(categoryDTOs, "Categories fetched successfully")
}
