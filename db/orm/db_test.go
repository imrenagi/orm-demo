package orm_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/imrenagi/orm-demo/db/orm"
	"github.com/imrenagi/orm-demo/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return sqlDB, gormDB, mock
}

func TestFindAllPerson(t *testing.T) {

	sqlDB, db, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(db)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id"}).
		AddRow(1, "Foo", "08123123", 1).
		AddRow(2, "Bar", "08918234", 1)

	mock.ExpectQuery("SELECT (.+) FROM \"persons\"").WillReturnRows(persons)

	data, err := r.FindAll(context.TODO())
	assert.Nil(t, err)
	assert.Len(t, data, 2)

	assert.Nil(t, mock.ExpectationsWereMet())

}

func TestFindAllWithPreload(t *testing.T) {

	t.Skip() //TODO remove this skip

	sqlDB, db, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(db)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number"}).
		AddRow(1, "Foo", "08123123").
		AddRow(2, "Bar", "08918234")

	mock.ExpectQuery("SELECT (.+) FROM \"persons\"").
		WillReturnRows(persons)

	mock.ExpectQuery("SELECT (.+) FROM \"addresses\" WHERE \"addresses\".\"person_id\" IN").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "person_id"}).
			AddRow(1, "1600 villa st", 1).
			AddRow(2, "1601 villa st", 1).
			AddRow(3, "1600 villa st", 2).
			AddRow(4, "1601 villa st", 2))

	data, err := r.FindAll(context.TODO())

	assert.Nil(t, err)
	assert.Len(t, data, 2)

	phone1 := "08123123"
	phone2 := "08918234"

	expected := []models.Person{
		{
			ID:          1,
			Name:        "Foo",
			PhoneNumber: &phone1,
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
		},
		{
			ID:          2,
			Name:        "Bar",
			PhoneNumber: &phone2,
			Addresses: []models.Address{
				{
					ID:       3,
					Address:  "1600 villa st",
					PersonID: 2,
				},
				{
					ID:       4,
					Address:  "1601 villa st",
					PersonID: 2,
				},
			},
		},
	}

	assert.Equal(t, expected, data)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindByIDWithJoin(t *testing.T) {
	sqlDB, db, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(db)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id", "School__id", "School__name"}).
		AddRow(1, "Foo", "08123123", 1, 1, "tk merdeka")

	mock.ExpectQuery("SELECT (.+) FROM \"persons\" LEFT JOIN \"schools\" (.+) WHERE persons.id (.+) LIMIT 1").
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

func TestFindByIDWithCustomSelect(t *testing.T) {
	sqlDB, db, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(db)

	persons := sqlmock.NewRows([]string{"id", "name", "school_id", "School__id", "School__name"}).
		AddRow(1, "Foo", 1, 1, "tk merdeka")

	mock.ExpectQuery("SELECT persons.id, persons.name, persons.school_id, school.id, school.name FROM \"persons\" LEFT JOIN schools (.+) WHERE persons.id (.+) LIMIT 1").
		WithArgs(1).
		WillReturnRows(persons)

	data, err := r.FindByIDWithCustomSelect(context.TODO(), 1)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	expected := &models.Person{
		ID:       1,
		Name:     "Foo",
		SchoolID: 1,
		School: &models.School{
			ID:   1,
			Name: "tk merdeka",
		},
	}

	assert.Equal(t, expected, data)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindCompletedByID(t *testing.T) {
	sqlDB, db, mock := dbMock(t)
	defer sqlDB.Close()

	r := New(db)

	persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id", "School__id", "School__name", "PlaceOfBirth__ID", "PlaceOfBirth__City", "PlaceOfBirth__PersonID"}).
		AddRow(1, "Foo", "08123123", 1, 1, "tk merdeka", 1, "Jakarta", 1)

	addresses := sqlmock.NewRows([]string{"id", "address", "person_id"}).
		AddRow(1, "1600 villa st", 1).
		AddRow(2, "1601 villa st", 1)

	personGroups := sqlmock.NewRows([]string{"person_id", "group_id"}).
		AddRow(1, 1).
		AddRow(1, 2)

	groups := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "line dance").
		AddRow(2, "zumba dance")

	mock.ExpectQuery("SELECT (.+) FROM \"persons\" LEFT JOIN \"schools\" (.+) LEFT JOIN \"place_of_births\" (.+) WHERE persons.id (.+) LIMIT 1").
		WithArgs(1).
		WillReturnRows(persons)

	mock.ExpectQuery("SELECT (.+) FROM \"addresses\" WHERE \"addresses\".\"person_id\" (.+)").
		WithArgs(1).
		WillReturnRows(addresses)

	mock.ExpectQuery("SELECT (.+) FROM \"person_groups\" WHERE \"person_groups\".\"person_id\"").
		WithArgs(1).
		WillReturnRows(personGroups)

	mock.ExpectQuery("SELECT (.+) FROM \"groups\" WHERE \"groups\".\"id\" IN ").
		WithArgs(1, 2).
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
