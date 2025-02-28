package main

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/migrations"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/routes"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	config.Connect()

	migrations.Migrate()
	routes.RegisterAllRoutes(e)

	port := ":8080"
	log.Println("Server running on http://localhost:" + port)
	e.Logger.Fatal(e.Start(port))
}
