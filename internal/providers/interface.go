package providers

import "github.com/vkuzmenkova/currency-rates/models"

type RatesProvider interface {
	GetRate(base string, currency_code string) (*models.CurrencyRate, error)
}
