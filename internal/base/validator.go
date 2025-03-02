package base

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func BindAndValidate(c echo.Context, dto interface{}) (Response, bool) {
	if err := c.Bind(dto); err != nil {
		return SetErrorMessage("Invalid input", err.Error()), false
	}

	if err := c.Validate(dto); err != nil {
		validationErrors := formatValidationErrors(err)
		return SetErrorMessage("Validation error", validationErrors), false
	}

	return Response{}, true
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = "failed on the '" + fieldErr.Tag() + "' validation"
		}
	}
	return errors
}
