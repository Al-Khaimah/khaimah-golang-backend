package main

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"log"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/migrations"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
)

func main() {
	e := echo.New()

	e.Validator = &base.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	config.Connect()
	db := config.GetDB()
	migrations.Migrate()
	routes.RegisterAllRoutes(e, db)

	port := ":8080"
	log.Println("Server running on http://" + config.GetEnv("PUBLIC_IP", "0.0.0.0") + port)
	e.Logger.Fatal(e.Start("0.0.0.0" + port))
}
