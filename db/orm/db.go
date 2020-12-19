package orm

import (
	"context"

	"github.com/imrenagi/orm-demo/models"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

type DB struct {
	db *gorm.DB
}

func (d DB) FindAll(ctx context.Context) ([]models.Person, error) {

	var persons []models.Person
	tx := d.db.WithContext(ctx).Find(&persons)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return persons, nil
}
