package main

import (
	"log"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/migrations"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/routes"

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

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}

	migrations.Migrate()
	routes.RegisterAllRoutes(e, db)

	port := ":8080"
	log.Println("Server running on http://localhost:" + port)
	e.Logger.Fatal(e.Start(port))
}
