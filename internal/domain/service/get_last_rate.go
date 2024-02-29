package currencyrates

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (s *CurrenciesService) GetLastRate(ctx context.Context, base string, currencyCode string) (models.CurrencyRate, error) {
	sql, _, err := sq.Select("rate", "updated_at").From("currency_rates").
		Where(sq.Eq{"base": s.GetCode(base), "currency": s.GetCode(currencyCode)}).
		Where("updated_at is not null").
		OrderBy("updated_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Printf("sql select: %s\n", err)
		return models.CurrencyRate{}, fmt.Errorf("sql select: %w", err)
	}

	var updatedAt time.Time
	var value float64

	err = s.Repo.Conn.QueryRow(ctx, sql, s.GetCode(base), s.GetCode(currencyCode)).Scan(&value, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.CurrencyRate{}, NoValueFoundError{Currency: currencyCode}
	}
	if err != nil {
		log.Printf("QueryRow: %s", err)
		return models.CurrencyRate{}, fmt.Errorf("QueryRow: %w", err)
	}

	cr := models.CurrencyRate{Base: base, Currency: currencyCode, UpdatedAt: updatedAt.String(), Value: value}

	return cr, nil
}
