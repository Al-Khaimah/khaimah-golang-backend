package users

import (
	userHandler "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/handlers"
	userRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	userService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepo := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userService)

	authGroup := e.Group("/auth")
	// userGroup := e.Group("/user")

	authGroup.POST("/signup", userHandler.CreateUser)
}
