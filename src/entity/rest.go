package entity

import (
	"math"

	"github.com/irdaislakhuafa/go-sdk/operator"
)

type (
	HTTPMessage struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	MetaError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Meta struct {
		Path       string     `json:"path"`
		StatusCode int        `json:"statusCode"`
		StatusStr  string     `json:"statusStr"`
		Message    string     `json:"message"`
		Timestamp  string     `json:"timestamp"`
		Error      *MetaError `json:"error"`
		RequestID  string     `json:"requestId"`
	}

	Pagination struct {
		CurrentPage     int      `json:"currentPage"`
		CurrentElements int      `json:"currentElements"`
		TotalPages      int      `json:"totalPages"`
		TotalElements   int      `json:"totalElements"`
		SortBy          []string `json:"sortBy"`
		CursorStart     *string  `json:"cursorStart"`
		CursorEnd       *string  `json:"cursorEnd"`
	}

	PaginationParams struct {
		Limit     int    `json:"limit" query:"limit"`
		Page      int    `json:"page" query:"page"`
		OrderBy   string `json:"order_by" query:"order_by"`
		OrderType string `json:"order_type" query:"order_type"`
	}

	HTTPRes struct {
		Message    HTTPMessage `json:"message"`
		Meta       Meta        `json:"meta"`
		Data       any         `json:"data,omitempty"`
		Pagination *Pagination `json:"pagination,omitempty"`
	}
)

func GenPagination(page, limit, totalItems int, sortBy []string) Pagination {
	if limit <= 0 {
		limit = 15
	}
	if page < 0 {
		page = 0
	}

	if limit > totalItems {
		limit = totalItems
	}

	var totalPages float64
	if totalItems != 0 {
		totalPages = math.Ceil(float64(totalItems)/float64(limit)) - 1
	} else {
		totalPages = 0
	}
	return Pagination{
		CurrentPage:     page,
		CurrentElements: limit,
		TotalPages:      int(totalPages),
		TotalElements:   totalItems,
		SortBy:          operator.Ternary(len(sortBy) == 0, []string{"id DESC"}, sortBy),
		CursorStart:     new(string),
		CursorEnd:       new(string),
	}
}

func (p *PaginationParams) Parse() error {
	if p.Limit == 0 {
		p.Limit = 15
	}
	p.Page = (p.Page * p.Limit)

	if p.OrderBy == "" {
		p.OrderBy = "id"
	}

	if p.OrderType == "" {
		p.OrderType = "DESC"
	}

	return nil
}
