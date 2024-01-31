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
		WithMeta(map[string]interface{}{
			"dateFormat": "Y-m-d",
		}),
	)

	is.Equal("ID", field.Name)
	is.Equal("uuid", field.Attribute)
	is.Equal("text", field.Component)
	is.Equal("Y-m-d", field.Meta["dateFormat"])
	is.True(field.Sortable)
	is.True(field.Searchable)
	is.True(field.Visibility)
}
func TestActionField(t *testing.T) {
	is := assert.New(t)

	field := NewActionField("Filters", []*ActionItems{
		{
			Label: "Users",
			Link:  "/clients/{id}/users",
		},
		{
			Label: "Sites",
			Link:  "/clients/{id}/Sites",
		},
	})

	is.Equal("Users", field.Actions[0].Label)
	is.Equal("/clients/{id}/users", field.Actions[0].Link)
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
