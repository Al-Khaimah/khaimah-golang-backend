package base

import (
	"reflect"

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
		validationErrors := formatValidationErrors(err, dto)
		return SetErrorMessage("Validation error", validationErrors), false
	}

	return Response{}, true
}

/*
@NOTE: This function is used to format the validation errors.
It is uses the custom message from the struct tags.
If the custom message is not set, it will use the default message.
*/
func formatValidationErrors(err error, dto interface{}) string {
	var errors string
	var dt reflect.Type
	if dto != nil {
		dt = reflect.TypeOf(dto)
		if dt.Kind() == reflect.Ptr {
			dt = dt.Elem()
		}
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			if dto != nil {
				if field, found := dt.FieldByName(fieldErr.Field()); found {
					customMsg := field.Tag.Get("message")
					if customMsg != "" {
						errors = customMsg
						continue
					}
				}
			}
			errors = "failed on the '" + fieldErr.Tag() + "' validation"
		}
	}
	return errors
}
