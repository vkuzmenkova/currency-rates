package currencyrates

import (
	"context"
	"fmt"
	"time"

	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vkuzmenkova/currency-rates/internal/providers/vat"
)

const (
	TTL = 1 * time.Minute //uuid TTl
)

func (s *CurrenciesService) UpdateRate(ctx context.Context, base string, currencyCode string) (uuid.UUID, error) {
	err := s.checkInput(base, currencyCode)
	if err != nil {
		return uuid.Nil, err
	}

	value, err := s.KV.Get(ctx, fmt.Sprintf("%s_%s", currencyCode, base)).Result()
	if err != nil {
		uuidUpdate := uuid.New()
		// Save info about currencies and UUID of the update
		err = s.KV.Set(ctx, fmt.Sprintf("%s_%s", currencyCode, base), uuidUpdate.String(), TTL).Err()
		if err != nil {
			log.Error().Msgf("Failed to set key with TTL: %s", err.Error())
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
		log.Error().Msgf("Job failed. getRate: %s", err)
		return fmt.Errorf("getBaseRate: %w", err)
	}

	// Store value in DB
	err = s.Repo.InsertRate(ctx, uuidUpdate, s.CurrencyList.GetValueByCode(c.Base), s.CurrencyList.GetValueByCode(c.Currency), c.Value)
	if err != nil {
		log.Error().Msgf("Job failed. InsertRate: %s", err)
		return fmt.Errorf("getBaseRate: %w", err)
	}

	// Delete temporary info about update job
	err = s.KV.Del(ctx, fmt.Sprintf("%s_%s", currencyCode, base)).Err()
	if err != nil {
		log.Error().Msgf("Job failed. Failed to delete key: %s", err.Error())
		return fmt.Errorf("redis del: %s", err.Error())
	}
	log.Info().Msgf("Finished update %s, dur=%s\n", uuidUpdate, time.Since(startTime))

	return nil
}
