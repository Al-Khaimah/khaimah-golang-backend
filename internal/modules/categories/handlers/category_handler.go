package categories

import (
	base "github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categoryDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/dtos"
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
	response := h.CategoryService.GetAllCategories()
	return c.JSON(response.HTTPStatus, response)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var category categoryDTO.Category
	if res, ok := base.BindAndValidate(c, &category); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.CategoryService.CreateCategory(category)
	return c.JSON(response.HTTPStatus, response)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	categoryID := c.Param("id")
	var categoryDTO categoryDTO.Category
	if res, ok := base.BindAndValidate(c, &categoryDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.CategoryService.UpdateCategory(categoryID, categoryDTO)
	return c.JSON(response.HTTPStatus, response)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	categoryID := c.Param("id")

	response := h.CategoryService.DeleteCategory(categoryID)
	return c.JSON(response.HTTPStatus, response)
}
