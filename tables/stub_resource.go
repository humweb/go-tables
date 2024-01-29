package tables

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserResource struct {
	AbstractResource
}

func NewUserResource(db *gorm.DB, req *http.Request) *UserResource {
	r := &UserResource{
		AbstractResource{
			Table:           "users",
			Db:              db,
			Request:         req,
			HasGlobalSearch: true,
		},
	}

	r.Fields = r.GetFields()
	r.Filters = r.GetFilters()

	return r
}

func (u *UserResource) GetModel() string {
	return "users"
}

func (u *UserResource) GetFields() []*Field {
	return []*Field{
		NewField("ID", WithSortable()),
		NewField("First name", WithSortable(), WithVisibility()),
		NewField("Last name", WithSortable(), WithSearchable()),
		NewField("Email", WithSortable()),
		NewField("Username", WithSortable()),
		NewField("Last login", WithSortable()),
	}
}

func (u *UserResource) GetFilters() []*Filter {
	return []*Filter{
		NewFilter("ID"),
		NewFilter("Client ID"),
	}
}

func (u *UserResource) ApplyFilter(db *gorm.DB) {

	//if clientId := chi.URLParam(u.Request, "client"); clientId != "" {
	//	db.Where("client_id = ?", clientId)
	//}
	//
	//if siteId := chi.URLParam(u.Request, "site"); siteId != "" {
	//	db.Joins("inner join sites_users ON sites_users.user_id = users.id").Where("sites_users.site_id = ?", siteId)
	//}
}

func (u *UserResource) WithGlobalSearch(db *gorm.DB, val string) {

	if v, err := strconv.Atoi(val); err == nil {
		db.Where("id = ?", v)
	} else {
		val = "%" + val + "%"
		db.Where(
			db.Where(db.Where("first_name ilike ?", val).Or("last_name ilike ?", val).Or("email ilike ?", val)),
		)
	}
}
