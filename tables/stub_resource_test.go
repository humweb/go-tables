package tables

// Basic imports
import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/humweb/go-tables/testutils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"testing"
)

type ResourceTestSuite struct {
	suite.Suite
}

func (suite *ResourceTestSuite) TestDefaultRequest() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users"`)
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" ORDER BY id ASC LIMIT $1`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(25).
		WillReturnRows(users)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)
	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(25, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestFilters() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&search[global]=foo", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE (first_name ilike $1 OR last_name ilike $2 OR email ilike $3)`)
	mock.ExpectQuery(expectedCountSQL).
		WithArgs("%foo%", "%foo%", "%foo%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE (first_name ilike $1 OR last_name ilike $2 OR email ilike $3) ORDER BY id ASC LIMIT $4`)
	mock.ExpectQuery(expectedSQL).
		WithArgs("%foo%", "%foo%", "%foo%", 30).
		WillReturnRows(users)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)
	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestPreloads() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30", nil)
	res := NewUserResource(db, request)

	res.Preloads = []*Preload{
		{
			Name: "Client",
		},
	}
	client := sqlmock.
		NewRows([]string{"id", "title", "description"}).
		AddRow(1, "cli", "desc")
	users := sqlmock.
		NewRows([]string{"id", "client_id", "first_name", "last_name", "username", "password"}).
		AddRow(1, 1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users"`)
	mock.ExpectQuery(expectedCountSQL).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" ORDER BY id ASC LIMIT $1`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(30).
		WillReturnRows(users)

	expectedClientSQL := regexp.QuoteMeta(`SELECT * FROM "clients" WHERE "clients"."id" = $1`)
	mock.ExpectQuery(expectedClientSQL).
		WithArgs(1).
		WillReturnRows(client)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)
	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestPreloadsWithCondition() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30", nil)
	res := NewUserResource(db, request)

	res.Preloads = []*Preload{
		{
			Name: "Client",
			Extra: func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			},
		},
	}
	client := sqlmock.
		NewRows([]string{"id", "title", "description"}).
		AddRow(1, "cli", "desc")
	users := sqlmock.
		NewRows([]string{"id", "client_id", "first_name", "last_name", "username", "password"}).
		AddRow(1, 1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users"`)
	mock.ExpectQuery(expectedCountSQL).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" ORDER BY id ASC LIMIT $1`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(30).
		WillReturnRows(users)

	expectedClientSQL := regexp.QuoteMeta(`SELECT id, name FROM "clients" WHERE "clients"."id" = $1`)
	mock.ExpectQuery(expectedClientSQL).
		WithArgs(1).
		WillReturnRows(client)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)
	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestGlobalIntFilter() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&search[global]=1", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)
	mock.ExpectQuery(expectedCountSQL).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY id ASC LIMIT $2`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(1, 30).
		WillReturnRows(users)

	var aryUsers []UserPrivate

	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)
	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestApplySearch() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&search[last_name]=bar", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE last_name ILIKE $1`)
	mock.ExpectQuery(expectedCountSQL).
		WithArgs("%bar%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE last_name ILIKE $1 ORDER BY id ASC LIMIT $2`)
	mock.ExpectQuery(expectedSQL).
		WithArgs("%bar%", 30).
		WillReturnRows(users)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)
	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestApplyIntSearch() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&search[id]=1", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)
	mock.ExpectQuery(expectedCountSQL).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY id ASC LIMIT $2`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(1, 30).
		WillReturnRows(users)

	var aryUsers []UserPrivate
	_, _ = res.Paginate(res, aryUsers)

	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestFilterApply() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&filters[id]=1", nil)
	res := NewUserResource(db, request)

	res.DefaultPerPage = 100

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)
	mock.ExpectQuery(expectedCountSQL).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY id ASC LIMIT $2`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(1, 30).
		WillReturnRows(users)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)

	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestDefaultPerPage() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?filters[id]=1", nil)
	res := NewUserResource(db, request)

	res.DefaultPerPage = 100

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)
	mock.ExpectQuery(expectedCountSQL).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY id ASC LIMIT $2`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(1, 100).
		WillReturnRows(users)

	var aryUsers []UserPrivate
	resp, _ := res.Paginate(res, aryUsers)

	records := resp["records"].([]UserPrivate)

	suite.Equal(uint(1), records[0].ID)
	suite.Equal("foo", records[0].FirstName)
	suite.Equal("bar", records[0].LastName)
	suite.Equal("baz", records[0].Username)

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(100, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestFlagVisibility() {
	sqlDB, db, _ := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&hidden=first_name", nil)
	res := NewUserResource(db, request)

	suite.Equal(true, res.Fields[1].Visible)
	res.FlagVisibility()
	suite.Equal(false, res.Fields[1].Visible)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResourceTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceTestSuite))
}
