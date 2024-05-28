package tables

// Basic imports
import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/humweb/go-tables/testutils"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type SearchTestSuite struct {
	suite.Suite
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *SearchTestSuite) TestDefaultRequest() {
	sqlDB, db, mock := testutils.DBMock(suite.T())
	defer sqlDB.Close()

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE first_name ILIKE (.+)$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	expectedIntSQL := "^SELECT (.+) FROM \"users\" WHERE first_name = (.+)$"
	mock.ExpectQuery(expectedIntSQL).WillReturnRows(users)

	search := &Search{
		Label:   "First name",
		Field:   "first_name",
		Value:   "foo",
		Enabled: true,
	}

	// Test string search
	var res []map[string]interface{}
	d := db.Table("users")
	search.ApplySearch(d)
	d.Find(&res)

	// Test numeric search
	search.Value = "1"
	d = db.Table("users")
	search.ApplySearch(d)
	d.Find(&res)

	suite.Nil(mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}
