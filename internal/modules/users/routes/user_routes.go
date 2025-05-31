package users

import (
	"net/http"

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
	bookmarksRepo := userRepository.NewBookmarkRepository(db)
	newAuthService := userService.NewAuthService(userRepo)
	newUserService := userService.NewUserService(userRepo, authRepo, bookmarksRepo)
	newUserHandler := userHandler.NewUserHandler(newUserService)
	newAuthHandler := userHandler.NewAuthHandler(newAuthService)

	authGroup := e.Group("/auth")
	authGroup.POST("/signup", newUserHandler.CreateUser)
	authGroup.POST("/login", newUserHandler.LoginUser)
	authGroup.POST("/oauth", newAuthHandler.OAuthLogin)
	authGroup.POST("/logout", newUserHandler.LogoutUser, middlewares.AuthMiddleware(authRepo))

	userGroup := e.Group("/user", middlewares.AuthMiddleware(authRepo))
	userGroup.GET("/profile", newUserHandler.GetUserProfile)
	userGroup.PUT("/profile", newUserHandler.UpdateUserProfile)
	userGroup.DELETE("/profile/delete-my-account", newUserHandler.DeleteMyAccount)
	userGroup.PATCH("/profile/preferences", newUserHandler.UpdateUserPreferences)
	userGroup.PATCH("/profile/password", newUserHandler.ChangePassword)
	userGroup.GET("/bookmarks", newUserHandler.GetUserBookmarks)
	userGroup.POST("/bookmarks/:podcast_id", newUserHandler.ToggleBookmarkPodcast)
	userGroup.GET("/downloads", newUserHandler.GetDownloadedPodcasts)

	adminGroup := e.Group("/admin", middlewares.AdminMiddleware())
	adminGroup.POST("/mark-user-admin/:user_id", newUserHandler.MarkUserAsAdmin)
	adminGroup.GET("/all-users", newUserHandler.GetAllUsers)
	adminGroup.DELETE("/user/:id", newUserHandler.DeleteUser)

	monitor := func(c echo.Context) error {
		return c.String(http.StatusOK, "khaimah is live")
	}
	e.GET("/", monitor)
	e.HEAD("/", monitor)

}
