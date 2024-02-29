package providers

import "github.com/vkuzmenkova/currency-rates/models"

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=RatesProvider
type RatesProvider interface {
	GetRate(base string, currencyCode string) (*models.CurrencyRate, error)
}
