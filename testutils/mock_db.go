package testutils

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBMock Helps with db testing, It generates a mock instance
func DBMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqldb,
		DriverName:           "postgres",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}
	return sqldb, gormdb, mock
}
