package main

import (
	"flag"
	"log"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/migrations"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Seed the database with test data")
	clearFlag := flag.Bool("clear", false, "Clear the database before seeding")
	flag.Parse()

	e := echo.New()

	e.Validator = &base.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	config.Connect()
	db := config.GetDB()
	migrations.Migrate()

	if *clearFlag {
		migrations.ClearDatabase(db)
	}

	if *seedFlag {
		migrations.SeedDatabase(db)
	}

	routes.RegisterAllRoutes(e, db)

	port := ":8080"
	log.Println("Server running on http://" + config.GetEnv("PUBLIC_IP", "0.0.0.0") + port)
	e.Logger.Fatal(e.Start("0.0.0.0" + port))
}
