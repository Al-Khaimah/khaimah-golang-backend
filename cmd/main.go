package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/migrations"

	bookmarks "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/bookmarks/routes"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/routes"
	notifications "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/notifications/routes"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/routes"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/routes"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	config.Connect()
	migrations.Migrate()

	users.RegisterRoutes(e)
	podcasts.RegisterRoutes(e)
	notifications.RegisterRoutes(e)
	categories.RegisterRoutes(e)
	bookmarks.RegisterRoutes(e)

	port := ":8080"
	log.Println("Server running on http://localhost:" + port)
	e.Logger.Fatal(e.Start(port))
}
