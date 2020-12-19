package models

import (
	"time"
)

type Person struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `gorm:"type:varchar;not null"`
	PhoneNumber  *string       `gorm:"type:text"`
	SchoolID     uint          `gorm:"not null"`
	PlaceOfBirth *PlaceOfBirth //has one relationship
	School       *School       //belongs to relationship
	Addresses    []Address     `gorm:"foreignKey:PersonID"`     //has many relationship
	Group        []Group       `gorm:"many2many:person_groups"` //many-to-many relationship
}

type PlaceOfBirth struct {
	ID       uint   `gorm:"primaryKey"`
	City     string `gorm:"type:varchar;not null"`
	PersonID uint
	Date     time.Time
}

type Address struct {
	ID       uint   `gorm:"primaryKey"`
	Address  string `gorm:"type:varchar;not null"`
	PersonID uint
}

type School struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar;not null"`
}

type Group struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
