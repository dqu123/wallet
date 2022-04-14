package models

import (
	"time"

	"github.com/dqu123/wallet-backend/database"
	"gorm.io/gorm"
)

type VirtualCard struct {
	gorm.Model

	// Wallet internal fields
	ProviderName string `json:"providerName"`

	ExternalVirtualCard
}

// External Extend API Fields
type ExternalVirtualCard struct {
	ProviderID         string     `json:"providerId" gorm:"column:provider_id;unique"`
	Status             string     `json:"status"`
	DisplayName        string     `json:"displayName"`
	ExpirationUTC      *time.Time `json:"expirationUtc"`
	Currency           string     `json:"currency"`
	LimitCents         int        `json:"limitCents"`
	BalanceCents       int        `json:"balanceCents"`
	LifetimeSpentCents int        `json:"lifetimeSpentCents"`
	LastFour           string     `json:"last4"`
}

func (cc *VirtualCard) Upsert(context database.Context) {
	db := context.DB
	db.Create(cc)
}
