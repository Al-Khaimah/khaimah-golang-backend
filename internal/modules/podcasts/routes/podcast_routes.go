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
	authRepo := userRepository.NewAuthRepository(db)
	userRepo := userRepository.NewUserRepository(db)
	bookmarksRepo := userRepository.NewBookmarkRepository(db)

	userService := userService.NewUserService(userRepo, authRepo, bookmarksRepo)

	podcastRepo := podcastRepository.NewPodcastRepository(db)
	podcastService := podcastService.NewPodcastService(podcastRepo)
	podcastHandler := podcastHandler.NewPodcastHandler(podcastService, userService)

	podcastGroup := e.Group("/podcasts", middlewares.AuthMiddleware(authRepo))
	podcastGroup.GET("/", podcastHandler.GetAllPodcasts)
	podcastGroup.GET("/recommended", podcastHandler.GetRecommendedPodcasts)
	podcastGroup.GET("/:id", podcastHandler.GetPodcastDetails)
	podcastGroup.POST("/:podcast_id/like", podcastHandler.LikePodcast)
	podcastGroup.GET("/category/:category_id/", podcastHandler.GetPodcastsByCategory)
	podcastGroup.POST("/:podcast_id/download", podcastHandler.DownloadPodcast)

}
