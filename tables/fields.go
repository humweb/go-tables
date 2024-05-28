package tables

import (
	"github.com/humweb/go-tables/utils"
)

// A Field represents a table field
type Field struct {
	Component    string                 `json:"component"`
	Attribute    string                 `json:"attribute"`
	Name         string                 `json:"name"`
	Sortable     bool                   `json:"sortable"`
	Searchable   bool                   `json:"searchable"`
	Visibility   bool                   `json:"visibility"`
	Visible      bool                   `json:"visible"`
	HasArraySort bool                   `json:"has_array_sort"`
	Actions      []*ActionItems         `json:"actions,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
}

type FieldOption func(*Field)

// NewField creates a new table field
func NewField(name string, opts ...FieldOption) *Field {
	s := &Field{
		Name:         name,
		Attribute:    utils.Slug(name),
		Component:    "text",
		Sortable:     false,
		Searchable:   false,
		Visibility:   false,
		HasArraySort: false,
		Visible:      true,
		Actions:      nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func NewActionField(name string, actions []*ActionItems) *Field {
	s := &Field{
		Name:       name,
		Attribute:  utils.Slug(name),
		Component:  "action-field",
		Sortable:   false,
		Searchable: false,
		Visibility: false,
		Visible:    true,
		Actions:    actions,
	}
	return s
}

//
// Filter Options
//

// SetVisibility sets the fields visibility flag to show or hide a column
func (f *Field) SetVisibility(flag bool) {
	f.Visible = flag
}

// WithAttribute is a Field option to set the database field
func WithAttribute(name string) FieldOption {
	return func(s *Field) {
		s.Attribute = name
	}
}

// WithFieldComponent allows you to override the default component type (text)
func WithFieldComponent(name string) FieldOption {
	return func(s *Field) {
		s.Component = name
	}
}

// WithSortable is a Field option to allow the field to be sorted
func WithSortable() FieldOption {
	return func(s *Field) {
		s.Sortable = true
	}
}

// WithSearchable is a Field option to allow the field to be searched
func WithSearchable() FieldOption {
	return func(s *Field) {
		s.Searchable = true
	}
}

// WithVisibility is a Field option to allow the field's visibility to be toggled
func WithVisibility() FieldOption {
	return func(s *Field) {
		s.Visibility = true
	}
}

// WithMeta add extra meta information for special field types
func WithMeta(data map[string]interface{}) FieldOption {
	return func(s *Field) {
		s.Meta = data
	}
}

// WithArraySort tell the system to sort the results by slice not sql query
func WithArraySort() FieldOption {
	return func(s *Field) {
		s.HasArraySort = true
	}
}
