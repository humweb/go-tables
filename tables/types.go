package tables

import (
	"github.com/humweb/go-tables/utils"
	"net/url"
	"strconv"
	"strings"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id desc"
	}
	return p.Sort
}

type TableRequest struct {
	Page         int               `json:"page"`
	PerPage      int               `json:"perPage"`
	Sort         string            `json:"sort"`
	Search       map[string]string `json:"search"`
	Filters      map[string]string `json:"filters"`
	GlobalFilter Filter            `json:"global_filter"`
}

func (r *TableRequest) Fill(req *url.URL) {
	query := req.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	r.Page = utils.DefaultInt(page, 1)

	perPage, _ := strconv.Atoi(query.Get("perPage"))
	r.PerPage = utils.DefaultInt(perPage, 25)

	r.Sort = utils.DefaultSort(query.Get("sort"), "id")

	r.SetFilterAndSearch(req)
}

func (r *TableRequest) SetFilterAndSearch(query *url.URL) {

	filters := make(map[string]string, strings.Count(query.RawQuery, "filters"))
	search := make(map[string]string, strings.Count(query.RawQuery, "search"))

	for key, val := range query.Query() {
		if strings.HasPrefix(key, "filters") {
			key = strings.TrimPrefix(strings.TrimSuffix(key, "]"), "filters[")
			filters[key] = val[0]
		} else if strings.HasPrefix(key, "search") {
			key = strings.TrimPrefix(strings.TrimSuffix(key, "]"), "search[")
			search[key] = val[0]
		}
	}
	r.Filters = filters
	r.Search = search
}
