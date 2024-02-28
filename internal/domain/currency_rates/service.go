package currencyrates

import (
	"context"
	"fmt"

	kvredis "github.com/go-redis/redis/v8"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/vkuzmenkova/currency-rates/internal/db"
	"github.com/vkuzmenkova/currency-rates/internal/domain"
	"github.com/vkuzmenkova/currency-rates/models"
)

type CurrenciesServicer interface {
	GetLastRate(ctx context.Context, uuid uuid.UUID, base string, currency_code string) (models.CurrencyRate, error)
	UpdateRate(ctx context.Context, uuid uuid.UUID, base string, currency_code string) error
}

type CurrenciesService struct {
	registry     *db.CurrenciesRepo
	Enqueuer     *work.Enqueuer
	Pool         *work.WorkerPool
	CurrencyList *domain.CurrencyList
	KV           *kvredis.Client
}

func NewCurrenciesService(ctx context.Context) (*CurrenciesService, error) {
	db, err := db.NewCurrenciesRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.NewCurrenciesRepo: %w", err)
	}

	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
	enqueuer := work.NewEnqueuer("currency_rates", redisPool)
	pool := work.NewWorkerPool(ctx, 10, "currency_rates", redisPool)

	rdb := kvredis.NewClient(&kvredis.Options{
		Addr:     "localhost:6379", // Адрес Redis сервера
		Password: "",               // Пароль, если не установлен, оставьте пустым
		DB:       0,                // Номер базы данных Redis
	})

	service := CurrenciesService{
		registry:     db,
		Enqueuer:     enqueuer,
		Pool:         pool,
		CurrencyList: domain.NewCurrencyList(),
		KV:           rdb,
	}

	service.Pool.Job("update_currency_rate", service.UpdateRateJobV2)
	//вынести start в другое место??
	service.Pool.Start()

	return &service, nil
}

func (s *CurrenciesService) GetCode(code string) uint8 {
	return s.CurrencyList.AvailableCurrencies[code]
}

func (s *CurrenciesService) CreateJobNameCodeBase(base string, code string) string {
	return fmt.Sprintf("%s_%s", code, base)
}
