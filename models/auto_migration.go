package models

import (
	"fmt"

	"github.com/dqu123/wallet-backend/database"
)

func AutoMigrateModels(context *database.Context) {
	db := context.DB
	models := []interface{}{
		&Transaction{},
		&VirtualCard{},
		&AccessToken{},
	}
	fmt.Println("AutoMigrating models")
	for _, model := range models {
		fmt.Printf("Migrating %#v\n", model)
		db.AutoMigrate(model)
	}
}
