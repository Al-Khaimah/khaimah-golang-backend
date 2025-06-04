package routes

import (
	_ "github.com/Al-Khaimah/khaimah-golang-backend/docs"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/routes"
	notifications "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/notifications/routes"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/routes"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/routes"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterAllRoutes(e *echo.Echo, db *gorm.DB) {
	users.RegisterRoutes(e, db)
	categories.RegisterRoutes(e, db)
	podcasts.RegisterRoutes(e, db)
	notifications.RegisterRoutes(e)

	RegisterSwaggerRoutes(e)
}
