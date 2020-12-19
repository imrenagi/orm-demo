package orm

import (
	"context"
	"database/sql"
	"log"

	"github.com/imrenagi/orm-demo/models"
)

func New(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

type DB struct {
	db *sql.DB
}

func (d DB) FindAll(ctx context.Context) ([]models.Person, error) {

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM people")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var persons []models.Person

	for rows.Next() {
		var person models.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.PhoneNumber, &person.SchoolID); err != nil {
			log.Fatal(err)
		}
		persons = append(persons, person)
	}
	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}
