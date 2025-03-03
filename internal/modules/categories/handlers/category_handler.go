package categories

import (
	categoryService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/services"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	CategoryService *categoryService.CategoryService
}

func NewCategoryHandler(categoryService *categoryService.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: categoryService}
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	categories := h.CategoryService.GetCategories()
	return c.JSON(categories.HTTPStatus, categories)
}
