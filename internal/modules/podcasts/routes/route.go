package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/middlewares"
	podcastHandler "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/handlers"
	podcastRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/repositories"
	podcastService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/services"
	userRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	userService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	authRepository := userRepository.NewAuthRepository(db)
	userRepository := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepository, authRepository)
	podcastRepository := podcastRepository.NewPodcastRepository(db)
	podcastService := podcastService.NewPodcastService(podcastRepository)
	podcastHandler := podcastHandler.NewPodcastHandler(podcastService, userService)

	podcastGroup := e.Group("/podcasts", middlewares.AuthMiddleware(authRepository))
	podcastGroup.GET("/", podcastHandler.GetAllPodcasts)
	podcastGroup.GET("/recommended", podcastHandler.GetRecommendedPodcasts)
}
