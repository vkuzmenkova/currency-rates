package currencyrates

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (s *CurrenciesService) GetLastRate(ctx context.Context, base string, currency_code string) (models.CurrencyRate, error) {
	sql, _, err := sq.Select("rate", "updated_at").From("currency_rates").
		Where(sq.Eq{"base": s.GetCode(base), "currency": s.GetCode(currency_code)}).
		Where("updated_at is not null").
		OrderBy("updated_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Fatalf("sql select: %s", err)
		return models.CurrencyRate{}, fmt.Errorf("sql select: %w", err)
	}

	var updatedAt time.Time
	var value float64
	err = s.registry.Conn.QueryRow(ctx, sql, s.GetCode(base), s.GetCode(currency_code)).Scan(&value, &updatedAt)
	if err == pgx.ErrNoRows {
		return models.CurrencyRate{}, NoValueFoundError{Currency: currency_code}
	}
	if err != nil {
		log.Fatalf("QueryRow: %s", err)
		return models.CurrencyRate{}, fmt.Errorf("QueryRow: %w", err)
	}
	cr := models.CurrencyRate{Base: base, Currency: currency_code, UpdatedAt: updatedAt.String(), Value: value}

	return cr, nil
}
