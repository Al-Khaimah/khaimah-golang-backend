package categories

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	middlewares "github.com/Al-Khaimah/khaimah-golang-backend/internal/middlewares"
	categoryHandler "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/handlers"
	categoryRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"
	categoryService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/services"
	authRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	newCategoryRepository := categoryRepository.NewCategoryRepository(db)
	newAuthRepository := authRepository.NewAuthRepository(db)
	newCategoryService := categoryService.NewCategoryService(newCategoryRepository)
	newCategoryHandler := categoryHandler.NewCategoryHandler(newCategoryService)

	categoryGroup := e.Group("/categories", middlewares.AuthMiddleware(newAuthRepository))
	categoryGroup.GET("/", newCategoryHandler.GetCategories)
	categoryGroup.POST("/", newCategoryHandler.CreateCategory)
	categoryGroup.PUT("/:id", newCategoryHandler.UpdateCategory)
	categoryGroup.DELETE("/:id", newCategoryHandler.DeleteCategory)
}
