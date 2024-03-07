package db

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (r *CurrenciesRepo) GetLastRate(ctx context.Context, base uint8, currencyCode uint8) (models.CurrencyRate, error) {
	sql, _, err := sq.Select("rate", "updated_at").From("currency_rates").
		Where(sq.Eq{"base": base, "currency": currencyCode}).
		Where("updated_at is not null").
		OrderBy("updated_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("sql select: %w", err)
	}

	var updatedAt time.Time
	var value float64

	err = r.Conn.QueryRow(ctx, sql, base, currencyCode).Scan(&value, &updatedAt)
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("QueryRow: %w", err)
	}

	cr := models.CurrencyRate{UpdatedAt: updatedAt.String(), Value: value}

	return cr, nil
}
