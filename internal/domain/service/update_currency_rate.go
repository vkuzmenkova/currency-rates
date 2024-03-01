package currencyrates

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vkuzmenkova/currency-rates/internal/providers/vat"
)

const (
	TTL = 10 * time.Minute //~2 min max per job x 5 retries //uuid TTl
)

func (s *CurrenciesService) UpdateRate(ctx context.Context, base string, currencyCode string) (uuid.UUID, error) {
	value, err := s.KV.Get(ctx, fmt.Sprintf("%s_%s", currencyCode, base)).Result()
	if err != nil {
		uuidUpdate := uuid.New()
		// Save info about currencies and UUID of the update
		err = s.KV.Set(ctx, fmt.Sprintf("%s_%s", currencyCode, base), uuidUpdate.String(), TTL).Err()
		if err != nil {
			log.Error().Msgf("Failed to set key with TTL:", err)
			return uuid.Nil, nil
		}

		// Initiate update
		_, err = s.Enqueuer.Enqueue("update_currency_rate", work.Q{"base": base, "currency_code": currencyCode, "uuid": uuidUpdate})
		if err != nil {
			return uuid.Nil, fmt.Errorf("EnqueueUnique: %w", err)
		}

		return uuidUpdate, nil
	}

	uuidValue, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, nil
	}
	return uuidValue, nil
}

func (s *CurrenciesService) UpdateRateJob(job *work.Job) error {
	uuidUpdate := job.ArgString("uuid")
	base := job.ArgString("base")
	currencyCode := job.ArgString("currency_code")

	ctx := context.Background()

	log.Info().Msgf("Start update %s\n", uuidUpdate)
	startTime := time.Now()

	// Get info from the provider
	provider := vat.NewVATProvider()
	c, err := provider.GetRate(base, currencyCode)
	if err != nil {
		log.Error().Msgf("Job failed. getBaseRate: %s", err)
		return fmt.Errorf("getBaseRate: %w", err)
	}

	// Store value in DB
	sql, _, err := sq.Insert("currency_rates").
		Columns("uuid", "base", "currency", "rate", "created_at", "updated_at").
		Values(uuidUpdate, s.CurrencyList.GetValueByCode(base), s.CurrencyList.GetValueByCode(currencyCode), c.Value, time.Now(), time.Now()).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error().Msgf("Job failed. sql insert request: %s", err)
		return fmt.Errorf("sql insert request: %w", err)
	}

	_, err = s.Repo.Conn.Exec(ctx, sql, uuidUpdate, s.CurrencyList.GetValueByCode(base), s.CurrencyList.GetValueByCode(currencyCode), c.Value, time.Now(), time.Now())
	if err != nil {
		log.Error().Msgf("Job failed. exec: %s", err)
		return fmt.Errorf("exec: %w", err)
	}

	// Delete temporary info about update job
	err = s.KV.Del(ctx, fmt.Sprintf("%s_%s", currencyCode, base)).Err()
	if err != nil {
		log.Error().Msgf("Job failed. Failed to delete key: %w", fmt.Sprintf("%s_%s", currencyCode, base), err)
		return fmt.Errorf("redis del: %w", err)
	}
	log.Info().Msgf("Finished update %s, dur=%s\n", uuidUpdate, time.Since(startTime))

	return nil
}
