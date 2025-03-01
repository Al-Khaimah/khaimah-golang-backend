package base

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BindAndValidate(c echo.Context, dto interface{}) error {
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return nil
}
