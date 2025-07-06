package gorm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
