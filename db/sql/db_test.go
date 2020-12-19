package orm_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/imrenagi/orm-demo/db/sql"
	"github.com/stretchr/testify/assert"
)

func dbMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return sqlDB, mock
}

func TestFindAllPerson(t *testing.T) {

	sqlDB, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(sqlDB)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id"}).
		AddRow(1, "Foo", "08123123", 1).
		AddRow(2, "Bar", "08918234", 1)

	mock.ExpectQuery("SELECT (.+) FROM people").WillReturnRows(persons)

	data, err := r.FindAll(context.TODO())
	assert.Nil(t, err)
	assert.Len(t, data, 2)

	assert.Nil(t, mock.ExpectationsWereMet())
}