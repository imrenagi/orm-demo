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

// FindAll List all rows in table but do not joins or preload the other relations
func (d DB) FindAll(ctx context.Context) ([]models.Person, error) {

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM persons")
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

// FindByIDWithJoin return 1 row but it tries to join the data whose "belongs to" relationship
func (d DB) FindByIDWithJoin(ctx context.Context, ID uint) (*models.Person, error) {

	var person models.Person
	var school models.School
	err := d.db.QueryRowContext(ctx, `
		SELECT p.*, s.* FROM persons p
		LEFT JOIN schools s on s.id = p.school_id
		WHERE id=?	
	`, ID).
		Scan(
			&person.ID, &person.Name, &person.PhoneNumber, &person.SchoolID,
			&school.ID, &school.Name,
		)
	if err != nil {
		return nil, err
	}
	person.School = &school
	return &person, nil
}

// FindCompletePersonByID return 1 row with all properties eagerly loaded and joined.
func (d DB) FindCompletePersonByID(ctx context.Context, ID uint) (*models.Person, error) {

	var person models.Person
	var school models.School
	var pob models.PlaceOfBirth
	err := d.db.QueryRowContext(ctx, `
		SELECT p.*, s.*, pob.id, pob.city, pob.person_id FROM persons p
		LEFT JOIN schools s on s.id = p.school_id
		LEFT JOIN place_of_births pob on p.id = pob.person_id
		WHERE id=?	
	`, ID).
		Scan(
			&person.ID, &person.Name, &person.PhoneNumber, &person.SchoolID,
			&school.ID, &school.Name,
			&pob.ID, &pob.City, &pob.PersonID,
		)
	if err != nil {
		return nil, err
	}
	person.School = &school
	person.PlaceOfBirth = &pob

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM addresses WHERE person_id =?`", person.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(&address.ID, &address.Address, &address.PersonID); err != nil {
			log.Fatal(err)
		}
		addresses = append(addresses, address)
	}
	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}
	person.Addresses = addresses

	rows2, err := d.db.QueryContext(ctx, `
		SELECT g.id, g.name FROM person_groups pg 
		LEFT JOIN groups g on g.id = pg.group_id
		WHERE pg.person_id =?`, person.ID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var groups []models.Group
	for rows2.Next() {
		var group models.Group
		if err := rows2.Scan(&group.ID, &group.Name); err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	rerr = rows2.Close()
	if rerr != nil {
		return nil, err
	}
	person.Groups = groups

	return &person, nil
}
