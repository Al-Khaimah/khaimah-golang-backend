package users

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/middlewares"
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
	authGroup.POST("/logout", newUserHandler.LogoutUser, middlewares.AuthMiddleware(authRepo))

	userGroup := e.Group("/user", middlewares.AuthMiddleware(authRepo))
	userGroup.GET("/profile", newUserHandler.GetUserProfile)
	userGroup.PUT("/profile", newUserHandler.UpdateUserProfile)
	userGroup.PATCH("/profile/password", newUserHandler.ChangePassword)
	userGroup.GET("", newUserHandler.GetAllUsers)
	userGroup.DELETE("/:id", newUserHandler.DeleteUser)
}
