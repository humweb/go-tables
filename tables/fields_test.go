package tables

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFieldDefaults(t *testing.T) {
	is := assert.New(t)
	field := NewField("ID")

	// assert equality
	is.Equal("ID", field.Name)
	is.Equal("id", field.Attribute)
	is.Equal("text", field.Component)
	is.False(field.Sortable)
	is.False(field.Searchable)
	is.False(field.Visibility)
}

func TestFieldOptions(t *testing.T) {
	is := assert.New(t)

	field := NewField(
		"ID",
		WithSortable(),
		WithAttribute("uuid"),
		WithSearchable(),
		WithVisibility(),
	)

	is.Equal("ID", field.Name)
	is.Equal("uuid", field.Attribute)
	is.Equal("text", field.Component)
	is.True(field.Sortable)
	is.True(field.Searchable)
	is.True(field.Visibility)
}

func TestUserResource(t *testing.T) {
	is := assert.New(t)

	resource := &UserResource{}
	fields := resource.GetFields()

	is.Equal("ID", fields[0].Name)
	is.Equal("First name", fields[1].Name)
	is.Equal("first_name", fields[1].Attribute)
	is.True(fields[1].Sortable)
	is.True(fields[2].Sortable)
	is.True(fields[3].Sortable)
	is.True(fields[4].Sortable)
	is.True(fields[5].Sortable)
}
