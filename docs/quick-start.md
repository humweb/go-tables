---
outline: deep
---

# Quick start

```bash
go get -u github.com/humweb/go-tables
```

## Setup Resource
This page demonstrates setting up a table resource.

```go
package resources

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserResource struct {
	tables.AbstractResource
}

// NewUserResource creates new table resource
func NewUserResource(db *gorm.DB, req *http.Request) *UserResource {
    r := &UserResource{
        tables.AbstractResource{
            Db:              db,
            Request:         req,
            HasGlobalSearch: true,
        },
    }

    r.Fields = r.GetFields()
    r.Filters = r.GetFilters()
    
    return r
}

// GetFields setup fields for table
func (u *UserResource) GetFields() []*tables.Field {
    return []*tables.Field{
        NewField("ID", tables.WithSortable()),
        NewField("First name", tables.WithSortable(), tables.WithVisibility()),
        NewField("Last name", tables.WithSortable(), tables.WithSearchable()),
        NewField("Email", tables.WithSortable()),
        NewField("Username", tables.WithSortable()),
        NewField("Last login", tables.WithSortable()),
    }
}

// GetFilters allows you to setup options filters
func (u *UserResource) GetFilters() []*tables.Filter {
    return []*tables.Filter{
        tables.NewFilter("ID"),
        tables.NewFilter("Client ID"),
    }
}

// ApplyFilters allows you to add custom filters based on things like route parameters
func (u *UserResource) ApplyFilter(db *gorm.DB) {

    if clientId := chi.URLParam(u.Request, "client"); clientId != "" {
        db.Where("client_id = ?", clientId)
    }
    
    if siteId := chi.URLParam(u.Request, "site"); siteId != "" {
        db.Joins("inner join sites_users ON sites_users.user_id = users.id").Where("sites_users.site_id = ?", siteId)
    }
}

// WithGlobalSearch setup query for global search
func (u *UserResource) WithGlobalSearch(db *gorm.DB, val string) {

    if v, err := strconv.Atoi(val); err == nil {
        db.Where("id = ?", v)
    } else {
        val = "%" + val + "%"
        db.Where(
            db.Where(db.Where("first_name ilike ?", val).
                Or("last_name ilike ?", val).
                Or("email ilike ?", val)),
        )
    }
}
```

## HTTP Handler
```go
func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

    resource := resources.NewUserResource(h.DB, r)
	
	// Preload relationships
    resource.Preloads = []*tables.Preload{{
        Name: "Owner",
        Extra: func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "email")
        }},
    }
	
	// Pass model and get results
	var clients []models.Client
    response, _ := resource.Paginate(resource, clients)
	
    _ = h.Inertia.Render(w, r, "Users", response)
}

```