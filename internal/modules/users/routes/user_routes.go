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
	authRepo := userRepository.NewAuthRepository(db)
	newUserService := userService.NewUserService(userRepo, authRepo)
	newUserHandler := userHandler.NewUserHandler(newUserService)

	authGroup := e.Group("/auth")
	authGroup.POST("/signup", newUserHandler.CreateUser)
	authGroup.POST("/login", newUserHandler.LoginUser)

	// userGroup := e.Group("/user")

}
