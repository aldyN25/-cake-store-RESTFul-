package common

import (
	"net/url"
	"strconv"
)

type PaginationResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PaginationRequest struct {
	Limit int
	Page  int
}

func (c *PaginationRequest) ParseQuery(v url.Values) {

	limit, page := v.Get("limit"), v.Get("page")
	if limit == "" {
		c.Limit = 10
	}

	if page == "" {
		c.Page = 1
	}

	if limit != "" {
		c.Limit, _ = strconv.Atoi(limit)
		if c.Limit <= 0 {
			c.Limit = 10
		}
	}

	if page != "" {
		c.Page, _ = strconv.Atoi(page)
		if c.Page <= 0 {
			c.Page = 1
		}
	}
}
