package orm_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	. "github.com/imrenagi/orm-demo/db/sql"
)

func BenchmarkFindAll_WithDB(b *testing.B) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s DB.name=%s password=%s sslmode=disable", "127.0.0.1", "5432", "users", "users", "users")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := New(db)

	for n := 0; n < b.N; n++ {
		r.FindAll(context.TODO())
	}
}

func BenchmarkFindByIDWithJoin_WithDB(b *testing.B) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s DB.name=%s password=%s sslmode=disable", "127.0.0.1", "5432", "users", "users", "users")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := New(db)

	for n := 0; n < b.N; n++ {
		r.FindByIDWithJoin(context.TODO(), 1)
	}
}

func BenchmarkFindCompletedByID_WithDB(b *testing.B) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s DB.name=%s password=%s sslmode=disable", "127.0.0.1", "5432", "users", "users", "users")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := New(db)

	for n := 0; n < b.N; n++ {
		r.FindCompletePersonByID(context.TODO(), 1)
	}
}
