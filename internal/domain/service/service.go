package currencyrates

import (
	"context"
	"fmt"

	"github.com/vkuzmenkova/currency-rates/configs"

	goredis "github.com/go-redis/redis/v8"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/vkuzmenkova/currency-rates/internal/db"
	"github.com/vkuzmenkova/currency-rates/internal/domain"
	"github.com/vkuzmenkova/currency-rates/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=Service
type Service interface {
	GetLastRate(ctx context.Context, uuid uuid.UUID, base string, currencyCode string) (models.CurrencyRate, error)
	UpdateRate(ctx context.Context, uuid uuid.UUID, base string, currencyCode string) error
	GetRateByUUID(ctx context.Context, uuid uuid.UUID) (models.CurrencyRate, error)
}

type CurrenciesService struct {
	Repo         *db.CurrenciesRepo
	Enqueuer     *work.Enqueuer
	Pool         *work.WorkerPool
	CurrencyList *domain.CurrencyList
	KV           *goredis.Client
}

func NewCurrenciesService(ctx context.Context, config configs.Config) (*CurrenciesService, error) {
	repo, err := db.NewCurrenciesRepo(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("db.NewCurrenciesRepo: %w", err)
	}

	redisPool := &redis.Pool{
		MaxActive: 10,
		MaxIdle:   10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port))
		},
	}
	enqueuer := work.NewEnqueuer(config.Redis.Namespace, redisPool)
	pool := work.NewWorkerPool(ctx, 10, config.Redis.Namespace, redisPool)

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: "",
		DB:       0,
	})

	currencies, err := repo.GetCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo.GetCurrencies: %w", err)
	}

	service := CurrenciesService{
		Repo:         repo,
		Enqueuer:     enqueuer,
		Pool:         pool,
		CurrencyList: domain.NewCurrencyList("USD", &currencies),
		KV:           rdb,
	}

	service.Pool.JobWithOptions("update_currency_rate", work.JobOptions{MaxFails: config.Redis.JobRetry}, service.UpdateRateJob)
	service.Pool.Start()

	return &service, nil
}

func (s *CurrenciesService) GetCode(code string) uint8 {
	return s.CurrencyList.AvailableCurrencies[code]
}

func (s *CurrenciesService) checkInput(base string, code string) error {
	if !s.CurrencyList.IsCurrencyAvailable(base) || !s.CurrencyList.IsCurrencyAvailable(code) {
		return UnavailableCurrencyError{CurrencyList: s.CurrencyList.GetCurrencyListUpper()}
	}
	if base == code {
		return BaseAndCodeAreEqualError{Message: fmt.Sprintf("%s:%s", base, code)}
	}

	return nil
}
