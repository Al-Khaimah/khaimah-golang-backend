package base

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginationRequest struct {
	Page    int `query:"page,omitempty" form:"page" validate:"omitempty,min=1" default:"1"`
	PerPage int `query:"per_page,omitempty" form:"per_page" validate:"omitempty,min=1,max=100" default:"10"`
}

/*
BindPaginationParams extracts and validates pagination parameters from the request
page and per_page query parameters are optional and default to 1 and 10 respectively
*/
func (p *PaginationRequest) BindPaginationParams(c echo.Context) {
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			p.Page = page
		}
	}

	if perPageStr := c.QueryParam("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil {
			p.PerPage = perPage
		}
	}

	p.Page = p.getPage()
	p.PerPage = p.getPerPage()
}

func (p *PaginationRequest) getPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *PaginationRequest) getPerPage() int {
	if p.PerPage < 1 {
		return 10
	}
	if p.PerPage > 100 {
		return 100
	}
	return p.PerPage
}
