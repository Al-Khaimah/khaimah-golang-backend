package main

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/middlewares"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/migrations"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/routes"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	base.RegisterValidator(e)
	middlewares.RegisterAllGlobalMiddlewares(e)

	config.Connect()
	db := config.GetDB()
	migrations.Migrate()
	migrations.SeedDatabase(db)

	routes.RegisterAllRoutes(e, db)

	startServer(e)
}
