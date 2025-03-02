package base

import (
	"log"
	"net/http"
)

type Response struct {
	HTTPStatus         int         `json:"-"` // will not be included in JSON response
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

func newResponse(httpStatus int, messageType, title string, data interface{}, errors interface{}, description ...string) Response {
	var desc string
	if len(description) > 0 {
		desc = description[0]
	}

	if messageType == ErrorStatus || errors != nil {
		log.Printf("❌ ERROR: %s - %v", title, errors)
	}

	return Response{
		HTTPStatus:         httpStatus,
		MessageType:        messageType,
		MessageTitle:       title,
		MessageDescription: desc,
		Data:               data,
		Errors:             errors,
	}
}

func SetData(data interface{}) Response {
	return newResponse(http.StatusOK, SuccessStatus, "", data, nil)
}

func SetSuccessMessage(title string, description ...string) Response {
	return newResponse(http.StatusOK, SuccessStatus, title, nil, nil, description...)
}

func SetErrorMessage(title string, errDetails interface{}) Response {
	return newResponse(http.StatusBadRequest, ErrorStatus, title, nil, errDetails)
}

func SetWarningMessage(title string, description ...string) Response {
	return newResponse(http.StatusConflict, WarningStatus, title, nil, nil, description...)
}
