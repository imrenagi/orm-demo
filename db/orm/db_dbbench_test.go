package orm_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/imrenagi/orm-demo/db/orm"
	"github.com/imrenagi/orm-demo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkFindAll_WithDB(b *testing.B) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s DB.name=%s password=%s sslmode=disable", "127.0.0.1", "5432", "users", "users", "users")
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})

	gormDB.AutoMigrate(&models.School{}, &models.Group{}, &models.Person{}, &models.PlaceOfBirth{}, &models.Address{})

	r := New(gormDB)

	for n := 0; n < b.N; n++ {
		r.FindAll(context.TODO())
	}
}

func BenchmarkFindByIDWithJoin_WithDB(b *testing.B) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s DB.name=%s password=%s sslmode=disable", "127.0.0.1", "5432", "users", "users", "users")
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})

	r := New(gormDB)

	for n := 0; n < b.N; n++ {
		r.FindByIDWithJoin(context.TODO(), 1)
	}
}

func BenchmarkFindCompletedByID_WithDB(b *testing.B) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s DB.name=%s password=%s sslmode=disable", "127.0.0.1", "5432", "users", "users", "users")
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})

	r := New(gormDB)

	for n := 0; n < b.N; n++ {
		r.FindCompletePersonByID(context.TODO(), 1)
	}
}
