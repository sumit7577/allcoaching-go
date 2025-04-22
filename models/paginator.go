package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

type Pagination struct {
	Offset int
	Limit  int
	query  orm.QuerySeter
	Count  int64
}

type PaginationSerializer struct {
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Count    int64       `json:"count"`
	Data     interface{} `json:"data"`
}

func (p *Pagination) Paginate() orm.QuerySeter {
	return p.query.Offset(p.Offset * p.Limit).Limit(p.Limit)
}

func (p *Pagination) CreatePagination(data interface{}) (*PaginationSerializer, error) {
	count, err := p.query.Count()
	if err != nil {
		return nil, err
	}
	p.Count = count

	totalPages := (count + int64(p.Limit) - 1) / int64(p.Limit)
	currentPage := p.Offset

	var next, previous *string

	if currentPage < int(totalPages) && currentPage+1 != int(totalPages) {
		nextURL := fmt.Sprintf("%s?page=%d", "http://127.0.0.8080", currentPage+1)
		next = &nextURL
	}

	if currentPage > 1 {
		prevURL := fmt.Sprintf("%s?page=%d", "http://127.0.0.8080", currentPage-1)
		previous = &prevURL
	}

	return &PaginationSerializer{
		Next:     next,
		Previous: previous,
		Count:    totalPages,
		Data:     data,
	}, nil
}
