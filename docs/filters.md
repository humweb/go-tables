---
outline: deep
---

# Table Filters


## Define filters

```go
func (u *UserResource) GetFilters() []*Filter {
	return []*Filter{
		NewFilter("ID"),
		NewFilter("Client ID"), // Filters client_id
	}
}
```


## Filter Options

`FilterOpt` is an optional function type to set filter attributes

```go
type FilterOpt func(*Filter)
```

### WithField
Allows you to override the default generated (snake_case) field name

```go
tables.WithField(name string)
```

### WithComponent
allows you to override the default component type (text)

```go
tables.WithComponent(name string)
```

### WithOptions
Allows you to set options for required filter types (select)

```go
tables.WithOptions(options ...FilterOptions)
```

#### FilterOptions defines filter options
```go
type FilterOptions struct {
    Label string `json:"label"`
    Value any    `json:"value"`
}
```