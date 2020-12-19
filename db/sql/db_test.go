package orm_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/imrenagi/orm-demo/db/sql"
	"github.com/imrenagi/orm-demo/models"
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

	mock.ExpectQuery("SELECT (.+) FROM persons").WillReturnRows(persons)

	data, err := r.FindAll(context.TODO())
	assert.Nil(t, err)
	assert.Len(t, data, 2)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindByIDWithJoin(t *testing.T) {
	sqlDB, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(sqlDB)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "School__id", "School__id", "School__name"}).
		AddRow(1, "Foo", "08123123", 1, 1, "tk merdeka")

	mock.ExpectQuery("SELECT p.*, s.* FROM persons p LEFT JOIN schools s on s.id = p.school_id WHERE id=?").
		WithArgs(1).
		WillReturnRows(persons)

	data, err := r.FindByIDWithJoin(context.TODO(), 1)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	phone := "08123123"
	expected := &models.Person{
		ID:          1,
		Name:        "Foo",
		PhoneNumber: &phone,
		SchoolID:    1,
		School: &models.School{
			ID:   1,
			Name: "tk merdeka",
		},
	}

	assert.Equal(t, expected, data)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindCompletedByID(t *testing.T) {
	sqlDB, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(sqlDB)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id", "School__id", "School__name", "PlaceOfBirth__ID", "PlaceOfBirth__City", "PlaceOfBirth__PersonID"}).
		AddRow(1, "Foo", "08123123", 1, 1, "tk merdeka", 1, "Jakarta", 1)

	addresses := sqlmock.NewRows([]string{"id", "address", "person_id"}).
		AddRow(1, "1600 villa st", 1).
		AddRow(2, "1601 villa st", 1)

	groups := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "line dance").
		AddRow(2, "zumba dance")

	mock.ExpectQuery("SELECT p.*, s.*, pob.id, pob.city, pob.person_id FROM persons p LEFT JOIN schools s on s.id = p.school_id LEFT JOIN place_of_births pob on p.id = pob.person_id WHERE id=?").
		WithArgs(1).
		WillReturnRows(persons)

	mock.ExpectQuery("SELECT (.+) FROM addresses WHERE person_id (.+)").
		WithArgs(1).
		WillReturnRows(addresses)

	mock.ExpectQuery("SELECT g.id, g.name FROM person_groups pg LEFT JOIN groups g on g.id = pg.group_id WHERE pg.person_id ").
		WithArgs(1).
		WillReturnRows(groups)

	data, err := r.FindCompletePersonByID(context.TODO(), 1)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	phone := "08123123"
	expected := &models.Person{
		ID:          1,
		Name:        "Foo",
		PhoneNumber: &phone,
		SchoolID:    1,
		School: &models.School{
			ID:   1,
			Name: "tk merdeka",
		},
		PlaceOfBirth: &models.PlaceOfBirth{
			ID:       1,
			City:     "Jakarta",
			PersonID: 1,
		},
		Addresses: []models.Address{
			{
				ID:       1,
				Address:  "1600 villa st",
				PersonID: 1,
			},
			{
				ID:       2,
				Address:  "1601 villa st",
				PersonID: 1,
			},
		},
		Groups: []models.Group{
			{
				ID:   1,
				Name: "line dance",
			},
			{
				ID:   2,
				Name: "zumba dance",
			},
		},
	}

	assert.Equal(t, expected, data)

	assert.Nil(t, mock.ExpectationsWereMet())
}
