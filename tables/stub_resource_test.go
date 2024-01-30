package tables

// Basic imports
import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/humweb/go-tables/testutils"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ResourceTestSuite struct {
	suite.Suite
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ResourceTestSuite) TestGetTablet() {
	sqlDB, db, _ := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	res := NewUserResource(db, request)

	suite.Equal("users", res.GetModel())
}

func (suite *ResourceTestSuite) TestDefaultRequest() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := "^SELECT (.+) FROM \"users\"$"
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := "^SELECT (.+) FROM \"users\" ORDER BY id ASC LIMIT 25$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	resp, _ := res.Paginate(res)

	records := resp["records"].([]map[string]interface{})
	suite.Equal(int64(1), records[0]["id"])
	suite.Equal("foo", records[0]["first_name"])
	suite.Equal("bar", records[0]["last_name"])
	suite.Equal("baz", records[0]["username"])
	suite.Equal("passwd", records[0]["password"])

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

	expectedCountSQL := "^SELECT (.+) FROM \"users\" WHERE \\(first_name ilike (.+) OR last_name ilike (.+) OR email ilike (.+)\\)$"
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE \\(first_name ilike (.+) OR last_name ilike (.+) OR email ilike (.+)\\) ORDER BY id ASC LIMIT 30$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	resp, _ := res.Paginate(res)

	records := resp["records"].([]map[string]interface{})
	suite.Equal(int64(1), records[0]["id"])
	suite.Equal("foo", records[0]["first_name"])
	suite.Equal("bar", records[0]["last_name"])
	suite.Equal("baz", records[0]["username"])
	suite.Equal("passwd", records[0]["password"])

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

	expectedCountSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+)$"
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+) ORDER BY id ASC LIMIT 30$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	resp, _ := res.Paginate(res)

	records := resp["records"].([]map[string]interface{})
	suite.Equal(int64(1), records[0]["id"])
	suite.Equal("foo", records[0]["first_name"])
	suite.Equal("bar", records[0]["last_name"])
	suite.Equal("baz", records[0]["username"])
	suite.Equal("passwd", records[0]["password"])

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

	expectedCountSQL := "^SELECT (.+) FROM \"users\" WHERE last_name ILIKE (.+)$"
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE last_name ILIKE (.+) ORDER BY id ASC LIMIT 30$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	resp, _ := res.Paginate(res)

	records := resp["records"].([]map[string]interface{})
	suite.Equal(int64(1), records[0]["id"])
	suite.Equal("foo", records[0]["first_name"])
	suite.Equal("bar", records[0]["last_name"])
	suite.Equal("baz", records[0]["username"])
	suite.Equal("passwd", records[0]["password"])

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

	expectedCountSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+)$"
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+) ORDER BY id ASC LIMIT 30$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	_, _ = res.Paginate(res)

	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestFilterApply() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&filters[id]=1", nil)
	res := NewUserResource(db, request)

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedCountSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+)$"
	mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+) ORDER BY id ASC LIMIT 30$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	resp, _ := res.Paginate(res)

	records := resp["records"].([]map[string]interface{})
	suite.Equal(int64(1), records[0]["id"])
	suite.Equal("foo", records[0]["first_name"])
	suite.Equal("bar", records[0]["last_name"])
	suite.Equal("baz", records[0]["username"])
	suite.Equal("passwd", records[0]["password"])

	pagination := resp["pagination"].(map[string]interface{})

	suite.Equal(30, pagination["perPage"])
	suite.Equal(1, pagination["page"])
	suite.Equal(1, pagination["total_pages"])
	suite.Equal(int64(1), pagination["record_count"])
	suite.Nil(mock.ExpectationsWereMet())
}

func (suite *ResourceTestSuite) TestFFlagVisibility() {
	sqlDB, db, _ := testutils.DBMock(suite.T())
	defer sqlDB.Close()
	request, _ := http.NewRequest(http.MethodGet, "/users?perPage=30&hidden=first_name", nil)
	res := NewUserResource(db, request)

	//users := sqlmock.
	//	NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
	//	AddRow(1, "foo", "bar", "baz", "passwd")
	//
	//expectedCountSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+)$"
	//mock.ExpectQuery(expectedCountSQL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	//
	//expectedSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+) ORDER BY id ASC LIMIT 30$"
	//mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	suite.Equal(true, res.Fields[1].Visible)
	res.FlagVisibility()
	suite.Equal(false, res.Fields[1].Visible)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResourceTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceTestSuite))
}
