package db

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	sq "github.com/Masterminds/squirrel"
)

func (r *CurrenciesRepo) InsertRate(ctx context.Context, uuidUpdate string, base uint8, currencyCode uint8, value float64) error {
	sql, _, err := sq.Insert("currency_rates").
		Columns("uuid", "base", "currency", "rate", "created_at", "updated_at").
		Values(uuidUpdate, base, currencyCode, value, time.Now(), time.Now()).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error().Msgf("Job failed. sql insert request: %s", err)
		return fmt.Errorf("sql insert request: %w", err)
	}

	_, err = r.Conn.Exec(ctx, sql, uuidUpdate, base, currencyCode, value, time.Now(), time.Now())
	if err != nil {
		log.Error().Msgf("Job failed. exec: %s", err)
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
