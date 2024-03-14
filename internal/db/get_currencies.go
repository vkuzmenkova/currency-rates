package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vkuzmenkova/currency-rates/models"
)

func (r *CurrenciesRepo) GetCurrencies(ctx context.Context) (map[string]uint8, error) {
	sql, _, err := sq.Select("currency_code", "id").From("currencies").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}

	currencies := make(map[string]uint8)
	var c models.Currency

	rows, err := r.Conn.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&c)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		currencies[c.Code] = c.ID
	}

	return currencies, nil
}
