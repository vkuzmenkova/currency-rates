package models

import "github.com/google/uuid"

type CurrencyRate struct {
	UpdatedAt string  `json:"updated_at"`
	Base      string  `json:"base"`
	Currency  string  `json:"currency"`
	Value     float64 `json:"value"`
}

type CurrencyUpdateUUID struct {
	Base     string
	Currency string
	UUID     uuid.UUID
}

type Currency struct {
	Code string
	ID   uint8
}
