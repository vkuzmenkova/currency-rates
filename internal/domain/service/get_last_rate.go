package currencyrates

import (
	"context"

	"github.com/vkuzmenkova/currency-rates/models"
)

func (s *CurrenciesService) GetLastRate(ctx context.Context, base string, currencyCode string) (models.CurrencyRate, error) {
	err := s.checkInput(base, currencyCode)
	if err != nil {
		return models.CurrencyRate{}, err
	}

	cr, err := s.Repo.GetLastRate(ctx, s.GetCode(base), s.GetCode(currencyCode))
	//if errors.As(err, &pgx.ErrNoRows) {
	//	return models.CurrencyRate{}, NoValueFoundError{Currency: currencyCode}
	//}
	if err != nil {
		return models.CurrencyRate{}, err
	}
	cr.Base = base
	cr.Currency = currencyCode

	return cr, nil
}
