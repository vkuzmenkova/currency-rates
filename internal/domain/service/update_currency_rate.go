package currencyrates

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/vkuzmenkova/currency-rates/internal/providers/vat"
)

func (s *CurrenciesService) UpdateRate(ctx context.Context, base string, currencyCode string) (uuid.UUID, error) {
	timeout := 2 * time.Minute

	value, err := s.KV.Get(ctx, fmt.Sprintf("%s_%s", currencyCode, base)).Result()
	if err != nil {
		uuidUpdate := uuid.New()
		// Фиксируем запуск обновления и значение uuid
		err = s.KV.Set(ctx, fmt.Sprintf("%s_%s", currencyCode, base), uuidUpdate.String(), timeout).Err()
		if err != nil {
			fmt.Println("Failed to set key with TTL:", err)
			return uuid.Nil, nil
		}

		// Запускаем джобу обновения
		_, err = s.Enqueuer.EnqueueUnique("update_currency_rate", work.Q{"base": base, "currency_code": currencyCode, "uuid": uuidUpdate})
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

	log.Printf("Start update %s\n", uuidUpdate)

	// Получаем значение из провайдера
	provider := vat.NewVATProvider()
	c, err := provider.GetRate(base, currencyCode)
	if err != nil {
		log.Printf("getBaseRate: %s", err)
		return fmt.Errorf("getBaseRate: %w", err)
	}

	// Записываем значение в базу
	sql, _, err := sq.Insert("currency_rates").
		Columns("uuid", "base", "currency", "rate", "created_at", "updated_at").
		Values(uuidUpdate, s.CurrencyList.GetValueByCode(base), s.CurrencyList.GetValueByCode(currencyCode), c.Value, time.Now(), time.Now()).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Printf("sql insert request: %s", err)
		return fmt.Errorf("sql insert request: %w", err)
	}

	_, err = s.Repo.Conn.Exec(ctx, sql, uuidUpdate, s.CurrencyList.GetValueByCode(base), s.CurrencyList.GetValueByCode(currencyCode), c.Value, time.Now(), time.Now())
	if err != nil {
		log.Printf("exec: %s", err)
		return fmt.Errorf("exec: %w", err)
	}

	// Удаляем значение обновления
	err = s.KV.Del(ctx, fmt.Sprintf("%s_%s", currencyCode, base)).Err()
	if err != nil {
		log.Print("Failed to delete key: %w", fmt.Sprintf("%s_%s", currencyCode, base))
		return fmt.Errorf("redis del: %w", err)
	}
	log.Printf("Finished update %s\n", uuidUpdate)

	return nil
}
