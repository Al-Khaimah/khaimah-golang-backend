package base

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

func SetData(data interface{}, title ...string) Response {
	var tit string
	if len(title) > 0 {
		tit = title[0]
	}

	return newResponse(http.StatusOK, SuccessStatus, tit, data, nil)
}

func SetSuccessMessage(title string, description ...string) Response {
	return newResponse(http.StatusOK, SuccessStatus, title, nil, nil, description...)
}

func SetErrorMessage(title string, errDetails ...interface{}) Response {
	return newResponse(http.StatusBadRequest, ErrorStatus, title, nil, errDetails)
}

func SetWarningMessage(title string, description ...string) Response {
	return newResponse(http.StatusConflict, WarningStatus, title, nil, nil, description...)
}

// @NOTE: Preferably use SetPaginatedResponse instead since it uses post-db retrieval pagination.
func SetDataPaginated(c echo.Context, data []interface{}) Response {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(c.QueryParam("perPage"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	totalRecords := len(data)
	totalPages := (totalRecords + perPage - 1) / perPage
	startIndex := (page - 1) * perPage
	endIndex := startIndex + perPage

	if startIndex > totalRecords {
		startIndex = totalRecords
	}
	if endIndex > totalRecords {
		endIndex = totalRecords
	}

	paginatedData := data[startIndex:endIndex]

	responseData := map[string]interface{}{
		"data":          paginatedData,
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     perPage,
	}

	return newResponse(http.StatusOK, SuccessStatus, "Success", responseData, nil)
}

func SetPaginatedResponse(data []interface{}, page, perPage, totalCount int) Response {
	totalPages := (totalCount + perPage - 1) / perPage

	responseData := map[string]interface{}{
		"items": data,
		"pagination": map[string]interface{}{
			"current_page": page,
			"page_size":    perPage,
			"total_items":  totalCount,
			"total_pages":  totalPages,
		},
	}
	return newResponse(http.StatusOK, SuccessStatus, "Success", responseData, nil)
}
