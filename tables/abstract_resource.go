package tables

import (
	"cmp"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/humweb/go-tables/utils"
	"gorm.io/gorm"
)

type AbstractResource struct {
	DB              *gorm.DB
	Request         *http.Request
	Fields          []*Field
	Filters         []*Filter
	Searches        []*Search
	Preloads        []Preload
	TableRequest    *TableRequest
	HasGlobalSearch bool
	DefaultPerPage  int
}

type Response struct {
	Records    any        `json:"records"`
	TableProps TableProps `json:"tableProps"`
	Pagination Pagination `json:"pagination"`
}

type TableProps struct {
	Sort    string             `json:"sort"`
	Page    int                `json:"page"`
	PerPage int                `json:"perPage"`
	Columns []*Field           `json:"columns"`
	Search  map[string]*Search `json:"search"`
	Filters []*Filter          `json:"filters"`
}

func (r *AbstractResource) ToResponse(paged *Pagination) Response {
	r.FlagVisibility()

	return Response{
		Records: paged.Rows,
		TableProps: TableProps{
			Sort:    utils.DefaultString(r.Request.URL.Query().Get("sort"), "id"),
			Page:    paged.Page,
			PerPage: paged.Limit,
			Columns: r.Fields,
			Search:  r.collectFieldSearches(),
			Filters: r.Filters,
		},
		Pagination: Pagination{
			Limit:      paged.Limit,
			Page:       paged.Page,
			TotalPages: paged.TotalPages,
			TotalRows:  paged.TotalRows,
		},
	}
}

// collectFieldSearches populates searches map from searchable fields and global search
func (r *AbstractResource) collectFieldSearches() map[string]*Search {
	var (
		searches = make(map[string]*Search)
		val      string
		ok       bool
	)

	// If global search enabled, we should always show it
	if r.HasGlobalSearch {
		val = r.TableRequest.Search["global"]
		searches["global"] = &Search{
			Label:   "Search..",
			Field:   "global",
			Value:   val,
			Enabled: true,
		}
	}

	// Handle Searchable fields
	for _, field := range r.Fields {
		if field.Searchable {
			val, ok = r.TableRequest.Search[field.Attribute]
			searches[field.Attribute] = &Search{
				Label:   field.Name,
				Field:   field.Attribute,
				Value:   val,
				Enabled: ok,
			}
		}
	}

	return searches
}

// FlagVisibility applies visibility flag to field the attributes
func (r *AbstractResource) FlagVisibility() {
	fields := r.Request.URL.Query().Get("hidden")
	fieldS := strings.Split(fields, ",")

	for k, val := range r.Fields {
		if slices.Contains(fieldS, val.Attribute) {
			r.Fields[k].SetVisibility(false)
		}
	}
}

// ApplySearch applies searchc criteria to query
func (r *AbstractResource) ApplySearch(db *gorm.DB, field, value string) {
	if v, err := strconv.Atoi(value); err == nil {
		db.Where(field+" = ?", v)
	} else {
		db.Where(field+" ILIKE ?", "%"+value+"%")
	}
}

// Paginate this is the main function for our resource
// It applies filters and search criteria and paginates
// Pagination uses a "Length aware" approach
func (r *AbstractResource) Paginate(resource ITable, model any) (Response, error) {
	r.TableRequest = &TableRequest{}

	var totalRows int64

	// Parse filters and search from request
	r.TableRequest.Fill(r.Request.URL)

	if r.TableRequest.PerPage == 25 && r.DefaultPerPage != 0 {
		r.TableRequest.PerPage = r.DefaultPerPage
	}

	// Init pagination
	p := &Pagination{
		Limit: r.TableRequest.PerPage,
		Page:  r.TableRequest.Page,
		Sort:  r.TableRequest.Sort,
	}

	// -- Start Query
	q := r.DB.Debug().Model(model)

	// Apply filters to query
	r.applySearch(resource, q)
	r.applyFilters(r.TableRequest.Filters, q)

	resource.ApplyFilter(q)

	// -- Get records count
	q.Count(&totalRows)
	p.TotalRows = totalRows

	// Eager load relationships
	r.eagerLoad(q)

	// Start pagination
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.Limit)))
	p.TotalPages = totalPages

	shouldArraySort := false
	arraySortField := ""
	for _, field := range r.Fields {
		if strings.Contains(p.GetSort(), field.Attribute) && field.HasArraySort {
			arraySortField = field.Attribute
			shouldArraySort = true
			if strings.Contains(p.GetSort(), "DESC") {
				arraySortField = "-" + arraySortField
			}
			p.Sort = ""
		}
	}

	// add pagination offset and order
	q.Offset(p.GetOffset()).
		Limit(p.GetLimit()).
		Order(p.GetSort())

	var data []map[string]any

	// Get results
	err := q.Debug().Find(&data).Error

	if err == nil {
		fmt.Println("BEFORE:", data)
		if shouldArraySort {
			r.sortArray(data, arraySortField)
		}
		fmt.Println("AFTER:", data)
		p.Rows = data
	}

	return r.ToResponse(p), err
}

func (r *AbstractResource) sortArray(data []map[string]any, field string) {
	if field[0:1] == "-" {
		field = field[1:]
		slices.SortStableFunc(data, func(a map[string]any, b map[string]any) int {
			return cmp.Compare(fmt.Sprintf("%v", b[field]), fmt.Sprintf("%v", a[field]))
		})
	} else {
		slices.SortStableFunc(data, func(a map[string]any, b map[string]any) int {
			return cmp.Compare(fmt.Sprintf("%v", a[field]), fmt.Sprintf("%v", b[field]))
		})
	}
}

// applyFilters applies filter criteria to the database query
func (r *AbstractResource) applyFilters(filters map[string]string, q *gorm.DB) {
	for _, f := range r.Filters {
		if val, ok := filters[f.Field]; ok {
			f.Value = val
			f.ApplyQuery(q)
		}
	}
}

// applySearch applies search criteria to the database query
func (r *AbstractResource) applySearch(resource ITable, q *gorm.DB) {

	for field, value := range r.TableRequest.Search {
		if field == "global" {
			resource.WithGlobalSearch(q, value)
		} else {
			r.ApplySearch(q, field, value)
		}
	}
}

// applySearch applies search criteria to the database query
func (r *AbstractResource) eagerLoad(q *gorm.DB) {
	for _, rel := range r.Preloads {
		if rel.Extra == nil {
			fmt.Println("preload", r.Preloads)

			q.Preload(rel.Name)
		} else {
			q.Preload(rel.Name, rel.Extra)
		}
	}
}
