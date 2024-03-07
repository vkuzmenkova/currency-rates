package db

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (r *CurrenciesRepo) GetRateByUUID(ctx context.Context, uuid uuid.UUID) (models.CurrencyRate, error) {
	sql, _, err := sq.Select("c_base.currency_code", "c_code.currency_code", "currency_rates.rate", "currency_rates.updated_at").From("currency_rates").
		Join("currencies as c_base on currency_rates.base=c_base.id").
		Join("currencies as c_code on currency_rates.currency=c_code.id").
		Where(sq.Eq{"uuid": uuid}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("sql select: %w", err)
	}

	var updatedAt time.Time
	var value float64
	var baseVal, currencyVal string

	err = r.Conn.QueryRow(ctx, sql, uuid).Scan(&baseVal, &currencyVal, &value, &updatedAt)
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("QueryRow: %w", err)
	}

	cr := models.CurrencyRate{UpdatedAt: updatedAt.String(), Value: value, Base: baseVal, Currency: currencyVal}

	return cr, nil
}
