package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterSwaggerRoutes(e *echo.Echo) {
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
