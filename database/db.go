package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}

// InitDB initializes the database connection
func InitDB() (*Context, error) {
	dsn := "postgresql://localhost/wallet"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	context := &Context{
		DB: db,
	}
	return context, nil
}
