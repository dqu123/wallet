package models

import (
	"time"

	"gorm.io/gorm"
)

type AccessToken struct {
	gorm.Model

	TokenValue    string
	ExpirationUTC *time.Time
}
