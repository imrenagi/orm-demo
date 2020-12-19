package orm

import (
	"context"

	"github.com/imrenagi/orm-demo/models"
	"gorm.io/gorm"
)

// New ...
func New(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// DB ...
type DB struct {
	db *gorm.DB
}

// FindAll List all rows in table but do not joins or preload the other relations
func (d DB) FindAll(ctx context.Context) ([]models.Person, error) {

	var persons []models.Person
	tx := d.db.WithContext(ctx).
		Find(&persons)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return persons, nil
}

// FindAllWithPreload list all rows in table, but it tries to preload an array to see
// whether eager preload is done per person basis or once for all.
func (d DB) FindAllWithPreload(ctx context.Context) ([]models.Person, error) {

	var persons []models.Person
	tx := d.db.WithContext(ctx).
		Preload("Addresses"). //TODO remove later
		Find(&persons)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return persons, nil
}

// FindByIDWithJoin return 1 row but it tries to join the data whose belongs to relationship
func (d DB) FindByIDWithJoin(ctx context.Context, ID uint) (*models.Person, error) {

	var person models.Person
	tx := d.db.WithContext(ctx).
		Where("persons.id = ?", ID).
		Joins("School").
		First(&person)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &person, nil
}

// FindByIDWithCustomSelect return 1 row but use select to pick fields that only matter
func (d DB) FindByIDWithCustomSelect(ctx context.Context, ID uint) (*models.Person, error) {

	var person models.Person
	tx := d.db.WithContext(ctx).
		Select("persons.id, persons.name, persons.school_id, school.id, school.name").
		Joins("LEFT JOIN schools on schools.id = persons.id").
		Where("persons.id = ?", ID).
		First(&person)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &person, nil
}

// FindCompletePersonByID return 1 row with all properties eagerly loaded and joined.
func (d DB) FindCompletePersonByID(ctx context.Context, ID uint) (*models.Person, error) {

	var person models.Person
	tx := d.db.WithContext(ctx).
		Where("persons.id = ?", ID).
		Joins("School").
		Joins("PlaceOfBirth").
		Preload("Addresses").
		Preload("Groups").
		First(&person)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &person, nil
}
