package models

import "github.com/google/uuid"

type CurrencyRate struct {
	UpdatedAt string
	Base      string
	Currency  string
	Value     float64
}

type CurrencyUpdateUUID struct {
	Base     string
	Currency string
	UUID     uuid.UUID
}
