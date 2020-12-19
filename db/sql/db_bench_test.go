package orm_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/imrenagi/orm-demo/db/sql"
)

func BenchmarkFindAll_WithMock(b *testing.B) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fail()
		b.Log(err)
	}
	defer sqlDB.Close()

	r := New(sqlDB)

	for n := 0; n < b.N; n++ {
		persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "school_id"}).
			AddRow(1, "Foo", "08123123", 1).
			AddRow(2, "Bar", "08918234", 1)

		mock.ExpectQuery("SELECT (.+) FROM persons").WillReturnRows(persons)

		r.FindAll(context.TODO())
	}
}

func BenchmarkFindByIDWithJoin_WithMock(b *testing.B) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fail()
		b.Log(err)
	}
	defer sqlDB.Close()

	r := New(sqlDB)

	for n := 0; n < b.N; n++ {
		persons := sqlmock.NewRows([]string{"id", "name", "phone_number", "School__id", "School__id", "School__name"}).
			AddRow(1, "Foo", "08123123", 1, 1, "tk merdeka")

		mock.ExpectQuery("SELECT p.*, s.* FROM persons p LEFT JOIN schools s on s.id = p.school_id WHERE id=?").
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
	defer sqlDB.Close()

	r := New(sqlDB)

	for n := 0; n < b.N; n++ {
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

		r.FindCompletePersonByID(context.TODO(), 1)
	}
}
