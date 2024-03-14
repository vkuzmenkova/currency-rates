package providers

import (
	"context"

	"github.com/vkuzmenkova/currency-rates/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=RatesProvider
type RatesProvider interface {
	GetRate(ctx context.Context, base string, currencyCode string) (*models.CurrencyRate, error)
}
