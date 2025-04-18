package categories

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	middlewares "github.com/Al-Khaimah/khaimah-golang-backend/internal/middlewares"
	categoryHandler "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/handlers"
	categoryRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"
	categoryService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/services"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	newCategoryRepository := categoryRepository.NewCategoryRepository(db)
	newCategoryService := categoryService.NewCategoryService(newCategoryRepository)
	newCategoryHandler := categoryHandler.NewCategoryHandler(newCategoryService)

	e.GET("/categories", newCategoryHandler.GetCategories)

	adminCategoryGroup := e.Group("/admin/categories", middlewares.AdminMiddleware())
	adminCategoryGroup.POST("/", newCategoryHandler.CreateCategory)
	adminCategoryGroup.PUT("/:id", newCategoryHandler.UpdateCategory)
	adminCategoryGroup.DELETE("/:id", newCategoryHandler.DeleteCategory)

}
