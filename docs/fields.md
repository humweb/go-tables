---
outline: deep
---

# Table Fields

Resource fields define the default table columns out put by the frontend.

## Define fields

```go
func (u *UserResource) GetFields() []*tables.Field {
	return []*tables.Field{
		tables.NewField(
			"ID",
			tables.WithSortable(),
        ),
		tables.NewField(
			"First name",
			tables.WithSortable(),
        ),
		tables.NewField(
			"Last name",
			tables.WithSortable(),
			tables.WithSearchable(),
        ),
		tables.NewField("Email",
			tables.WithSortable(),
        ),
		tables.NewField(
			"Last login",
			WithComponent("date"),
			tables.WithSortable(),
			tables.WithMeta(map[string]interface{}{
				"dateFormat": "MM/DD/YYYY ",
			}),
		),
	}
}
```

## Field Options
```go
type FieldOption func(*Field)
```

### `SetVisibility` 
Sets the fields visibility flag to show or hide a column
```go
tables.SetVisibility(bool)
```

### WithAttribute
Override attribute built from field title (`snake_case`)
```go
tables.WithAttribute(string)
```

### WithFieldComponent
Allows you to override the default component type (text)
```go
WithFieldComponent(string)
```

### WithSortable
Option to allow the field to be sorted
```go
tables.WithSortable
```

### WithSearchable
Option to allow the field to be searched
```go
tables.WithSearchable()
```

### WithVisibility
Option to allow the field's visibility to be toggled
```go
tables.WithVisibility(bool)
```

### WithMeta
Adds extra meta information for used by special field types, like the `date` field`
```go
tables.WithMeta(map[string]interface{})
```