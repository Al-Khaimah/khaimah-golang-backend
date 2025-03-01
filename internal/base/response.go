package base

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	MessageType        string      `json:"message_type,omitempty"`
	MessageTitle       string      `json:"message_title,omitempty"`
	MessageDescription string      `json:"message_description,omitempty"`
	Data               interface{} `json:"data,omitempty"`
	Errors             interface{} `json:"errors,omitempty"`
}

const (
	SuccessStatus = "success"
	WarningStatus = "warning"
	ErrorStatus   = "error"
)

func SetData(c echo.Context, data interface{}, httpStatus ...int) error {
	statusCode := http.StatusOK
	if len(httpStatus) > 0 {
		statusCode = httpStatus[0]
	}

	response := Response{
		Data: data,
	}
	return c.JSON(statusCode, response)
}

func SetSuccessMessage(c echo.Context, title string, description ...string) error {
	var desc string
	if len(description) > 0 {
		desc = description[0]
	} else {
		desc = ""
	}

	response := Response{
		MessageType:        SuccessStatus,
		MessageTitle:       title,
		MessageDescription: desc,
	}
	return c.JSON(http.StatusOK, response)
}

func SetErrorMessage(c echo.Context, title string, errDetails interface{}, httpStatus ...int) error {
	statusCode := http.StatusBadRequest
	if len(httpStatus) > 0 {
		statusCode = httpStatus[0]
	}

	response := Response{
		MessageType:  ErrorStatus,
		MessageTitle: title,
		Errors:       errDetails,
	}
	return c.JSON(statusCode, response)
}

func SetWarningMessage(c echo.Context, title string, description ...string) error {
	var desc string
	if len(description) > 0 {
		desc = description[0]
	} else {
		desc = ""
	}

	response := Response{
		MessageType:        WarningStatus,
		MessageTitle:       title,
		MessageDescription: desc,
	}
	return c.JSON(http.StatusBadRequest, response)
}
