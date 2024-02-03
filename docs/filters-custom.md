# Custom Filters
This resource hook allows you to set up custom filtering or any kind of query really.
Input could be based on url parameters or global variable

## ApplyFilter hook
```go
func (u *UserResource) ApplyFilter(db *gorm.DB) {

	if clientId := chi.URLParam(u.Request, "client"); clientId != "" {
		db.Where("client_id = ?", clientId)
	}

	if siteId := chi.URLParam(u.Request, "site"); siteId != "" {
		db.Joins("inner join sites_users ON site_users.user_id = users.id").
			Where("site_users.site_id = ?", siteId)
	}
}
```