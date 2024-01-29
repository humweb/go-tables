package tables

import (
	"strconv"

	"github.com/humweb/go-tables/utils"
	"gorm.io/gorm"
)

// Filter defines filters for adding query clauses to our query
type Filter struct {
	Component string          `json:"component"`
	Label     string          `json:"label"`
	Field     string          `json:"field"`
	Options   []FilterOptions `json:"options"`
	Value     string          `json:"value"`
}

// FilterOptions defines filter options
type FilterOptions struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// ApplyQuery adds search criteria to the database query
func (f *Filter) ApplyQuery(db *gorm.DB) {
	if v, err := strconv.Atoi(f.Value); err == nil {
		db.Where(f.Field+" = ?", v)
	} else {
		db.Where(f.Field+" ILIKE ?", "%"+f.Value+"%")
	}
}

// FilterOpt is an optional function type to set filter attributes
type FilterOpt func(*Filter)

// NewFilter creates a new filter
func NewFilter(name string, opts ...FilterOpt) *Filter {
	s := &Filter{
		Label:     name,
		Field:     utils.Slug(name),
		Component: "text",
		Value:     "",
		Options:   nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithField allows you to override the default generated field name
func WithField(name string) FilterOpt {
	return func(s *Filter) {
		s.Field = name
	}
}

// WithComponent allows you to override the default component type (text)
func WithComponent(name string) FilterOpt {
	return func(s *Filter) {
		s.Component = name
	}
}

// WithOptions allows you to set options for required filter types (select)
func WithOptions(options ...FilterOptions) FilterOpt {
	return func(s *Filter) {
		s.Options = options
	}
}
