package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	ProviderName string `json:"providerName"`
	ExternalTransaction
}

// Interface for external transactions
type ExternalTransaction struct {
	ProviderID             string     `json:"providerId" gorm:"column:provider_id;unique"`
	MCC                    string     `json:"mcc"`
	MCCGroup               string     `json:"mccGroup"`
	MCCDescription         string     `json:"mccDescription"`
	MerchantName           string     `json:"merchantName"`
	AuthBillingAmountCents int        `json:"authBillingAmountCents"`
	AuthBillingCurrency    string     `json:"authBillingCurrency"`
	AuthedUTC              *time.Time `json:"authedUtc"`
}
