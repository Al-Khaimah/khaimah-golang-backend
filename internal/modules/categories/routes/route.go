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
	categoryRepository := categoryRepository.NewCategoryRepository(db)
	authRepository := authRepository.NewAuthRepository(db)
	categoryService := categoryService.NewCategoryService(categoryRepository)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryService)

	categoryGroup := e.Group("/categories", middlewares.AuthMiddleware(authRepository))
	categoryGroup.GET("/", categoryHandler.GetCategories)
}
