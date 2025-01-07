package models

import "gorm.io/gorm"

type Books struct {
	// defined the id as the primary key and it will be generated by gorm itself
	ID        uint    `json:"id" gorm:"primaryKey; autoIncrement"`
	Author    *string `json:"author" gorm:"not null"`
	Title     *string `json:"title" gorm:"not null"`
	Publisher *string `json:"publisher" gorm:"not null"`
}

// auto migrate:  Automatically migrate your schema, to keep your schema up to date.
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}

// NOTE: AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
// It will change existing column’s type if its size, precision changed, or if it’s changing
// from non-nullable to nullable. It WON’T delete unused columns to protect your data.
