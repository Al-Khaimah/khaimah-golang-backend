package routes

import (
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/routes"
	notifications "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/notifications/routes"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/routes"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/routes"
	"github.com/labstack/echo/v4"
)

func RegisterAllRoutes(e *echo.Echo) {
	users.RegisterRoutes(e)
	podcasts.RegisterRoutes(e)
	notifications.RegisterRoutes(e)
	categories.RegisterRoutes(e)
}
