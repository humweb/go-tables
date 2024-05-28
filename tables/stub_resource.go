package tables

import (
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (Client) TableName() string {
	return "clients"
}

type UserPrivate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClientId  int       `json:"client_id"`
	Client    *Client   `gorm:"foreignkey:ClientId" json:"client,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserPrivate) TableName() string {
	return "users"
}

type UserResource struct {
	AbstractResource
}

func NewUserResource(db *gorm.DB, req *http.Request) *UserResource {
	r := &UserResource{
		AbstractResource{
			DB:              db,
			Request:         req,
			HasGlobalSearch: true,
		},
	}

	r.Fields = r.GetFields()
	r.Filters = r.GetFilters()

	return r
}

func (u *UserResource) GetFields() []*Field {
	return []*Field{
		NewField("ID", WithSortable()),
		NewField("First name", WithSortable(), WithVisibility(), WithArraySort()),
		NewField("Last name", WithSortable(), WithSearchable()),
		NewField("Email", WithSortable()),
		NewField("Username", WithSortable()),
		NewField("Last login", WithSortable()),
		NewActionField("Filters", []*ActionItems{
			{
				Label:  "Users",
				Link:   "/clients/{id}/users",
				Params: []string{"id"},
			},
			{
				Label:  "Sites",
				Link:   "/clients/{id}/Sites",
				Params: []string{"id"},
			},
		}),
	}
}

func (u *UserResource) GetFilters() []*Filter {
	return []*Filter{
		NewFilter("ID"),
		NewFilter("Client ID"),
	}
}

func (u *UserResource) ApplyFilter(db *gorm.DB) {
	// if clientId := chi.URLParam(u.Request, "client"); clientId != "" {
	//	db.Where("client_id = ?", clientId)
	//}
	//
	// if siteId := chi.URLParam(u.Request, "site"); siteId != "" {
	//	db.Joins("inner join sites_users ON sites_users.user_id = users.id").Where("sites_users.site_id = ?", siteId)
	//}
}

func (u *UserResource) WithGlobalSearch(db *gorm.DB, val string) {
	if v, err := strconv.Atoi(val); err == nil {
		db.Where("id = ?", v)
	} else {
		val = "%" + val + "%"
		db.Where("(first_name ilike ? OR last_name ilike ? OR email ilike ?)", val, val, val)
	}
}
