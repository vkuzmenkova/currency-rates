package currencyrates

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (s *CurrenciesService) GetRateByUUID(ctx context.Context, uuid uuid.UUID) (models.CurrencyRate, error) {
	cr, err := s.Repo.GetRateByUUID(ctx, uuid)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.CurrencyRate{}, NoUUIDFoundError{Message: err.Error()}
	}
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("QueryRow: %w", err)
	}

	return cr, nil
}
