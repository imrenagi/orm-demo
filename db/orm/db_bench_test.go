package orm_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/imrenagi/orm-demo/db/orm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkFindAll_WithMock(b *testing.B) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fail()
		b.Log(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		b.Fail()
		b.Log(err)
	}
	defer sqlDB.Close()

	r := New(gormDB)

	for n := 0; n < b.N; n++ {
		persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id"}).
			AddRow(1, "Foo", "08123123", 1).
			AddRow(2, "Bar", "08918234", 1)
		mock.ExpectQuery("SELECT (.+) FROM \"persons\"").WillReturnRows(persons)
		r.FindAll(context.TODO())
	}
}

func BenchmarkFindByIDWithJoin_WithMock(b *testing.B) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fail()
		b.Log(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		b.Fail()
		b.Log(err)
	}
	defer sqlDB.Close()

	r := New(gormDB)

	for n := 0; n < b.N; n++ {
		persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id", "School__id", "School__name"}).
			AddRow(1, "Foo", "08123123", 1, 1, "tk merdeka")

		mock.ExpectQuery("SELECT (.+) FROM \"persons\" LEFT JOIN \"schools\" (.+) WHERE persons.id (.+) LIMIT 1").
			WithArgs(1).
			WillReturnRows(persons)
		r.FindByIDWithJoin(context.TODO(), 1)
	}
}

func BenchmarkFindCompletedByID_WithMock(b *testing.B) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fail()
		b.Log(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		b.Fail()
		b.Log(err)
	}
	defer sqlDB.Close()

	r := New(gormDB)

	for n := 0; n < b.N; n++ {
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

		r.FindCompletePersonByID(context.TODO(), 1)
	}
}
