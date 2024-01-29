package tables

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/humweb/go-tables/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyQuery(t *testing.T) {
	is := assert.New(t)
	sqlDB, db, mock := testutils.DBMock(t)
	defer sqlDB.Close()

	f := &Filter{
		Component: "text",
		Label:     "Name",
		Field:     "first_name",
		Value:     "foo",
	}

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE first_name ILIKE (.+)$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	var user []map[string]interface{}

	d := db.Table("users")
	f.ApplyQuery(d)
	d.Find(&user)

	is.Nil(mock.ExpectationsWereMet())
}

func TestApplyQueryNumeric(t *testing.T) {
	is := assert.New(t)
	sqlDB, db, mock := testutils.DBMock(t)
	defer sqlDB.Close()

	f := &Filter{
		Component: "text",
		Label:     "ID",
		Field:     "id",
		Value:     "33",
	}

	users := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "username", "password"}).
		AddRow(1, "foo", "bar", "baz", "passwd")

	expectedSQL := "^SELECT (.+) FROM \"users\" WHERE id = (.+)$"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	var user []map[string]interface{}

	d := db.Table("users")
	f.ApplyQuery(d)
	d.Find(&user)

	is.Nil(mock.ExpectationsWereMet())
}

func TestNewFilter(t *testing.T) {
	is := assert.New(t)

	filter := NewFilter(
		"First name",
		WithField("foo"),
		WithComponent("select"),
		WithOptions(
			FilterOptions{Label: "Car", Value: 1},
			FilterOptions{Label: "Truck", Value: 2},
		),
	)

	is.Equal("First name", filter.Label)
	is.Equal("foo", filter.Field)
	is.Equal("select", filter.Component)
	is.Equal("Car", filter.Options[0].Label)
	is.Equal(1, filter.Options[0].Value)

	filterNoOptions := NewFilter("First name")
	is.Nil(filterNoOptions.Options)
}
