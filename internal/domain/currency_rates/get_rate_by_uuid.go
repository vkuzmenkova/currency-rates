package currencyrates

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (s *CurrenciesService) GetRateByUUID(ctx context.Context, uuid uuid.UUID) (models.CurrencyRate, error) {
	sql, _, err := sq.Select("base", "currency", "rate", "updated_at").From("currency_rates").
		Where(sq.Eq{"uuid": uuid}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Fatalf("sql select: %s", err)
		return models.CurrencyRate{}, fmt.Errorf("sql select: %w", err)
	}

	var updatedAt time.Time
	var value float64
	var baseVal, currencyVal uint8
	err = s.registry.Conn.QueryRow(ctx, sql, uuid).Scan(&baseVal, &currencyVal, &value, &updatedAt)
	if err == pgx.ErrNoRows {
		return models.CurrencyRate{}, NoUUIDFoundError{Message: err.Error()}
	}
	if err != nil {
		log.Fatalf("QueryRow: %s", err)
		return models.CurrencyRate{}, fmt.Errorf("QueryRow: %w", err)
	}
	cr := models.CurrencyRate{UpdatedAt: updatedAt.String(), Value: value, Base: s.CurrencyList.GetCodeByValue(baseVal), Currency: s.CurrencyList.GetCodeByValue(currencyVal)}

	return cr, nil
}
