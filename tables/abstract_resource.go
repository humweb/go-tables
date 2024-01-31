package tables

import (
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
	Model           Model
	Fields          []*Field
	Filters         []*Filter
	Searches        []*Search
	Preloads        []*Preload
	TableRequest    *TableRequest
	HasGlobalSearch bool
}

func (r *AbstractResource) ToResponse(paged *Pagination) map[string]interface{} {
	r.FlagVisibility()

	return map[string]interface{}{
		"records": paged.Rows,
		"tableProps": map[string]interface{}{
			"sort":    utils.DefaultString(r.Request.URL.Query().Get("sort"), "id"),
			"page":    paged.Page,
			"perPage": paged.Limit,
			"columns": r.Fields,
			"search":  r.collectFieldSearches(),
			"filters": r.Filters,
		},
		"pagination": map[string]interface{}{
			"perPage":      paged.Limit,
			"page":         paged.Page,
			"total_pages":  paged.TotalPages,
			"record_count": paged.TotalRows,
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
func (r *AbstractResource) Paginate(resource ITable) (map[string]interface{}, error) {
	r.TableRequest = &TableRequest{}

	var totalRows int64

	// Parse filters and search from request
	r.TableRequest.Fill(r.Request.URL)

	// Init pagination
	p := &Pagination{
		Limit: r.TableRequest.PerPage,
		Page:  r.TableRequest.Page,
		Sort:  r.TableRequest.Sort,
	}

	// -- Start Query
	q := r.DB.Table(r.Model.TableName())

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

	var res []map[string]interface{}
	q.Select(r.getSelectFields()).
		Offset(p.GetOffset()).
		Limit(p.GetLimit()).
		Order(p.GetSort())

	// Get results
	err := q.Debug().Find(&res).Error
	if err == nil {
		p.Rows = res
	}

	return r.ToResponse(p), err
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
		q.Preload(rel.Name, rel.Extra)
	}
}

// getSelectFields returns a select statement with only the fields added to the resource
func (r *AbstractResource) getSelectFields() string {
	ary := make([]string, len(r.Fields))
	for i, f := range r.Fields {
		if f.Component != "action-field" {
			ary[i] = r.Model.TableName() + "." + f.Attribute
		}
	}
	return strings.Trim(strings.Join(ary, ","), ",")
}
