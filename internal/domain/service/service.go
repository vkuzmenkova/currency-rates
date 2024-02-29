package currencyrates

import (
	"context"
	"fmt"

	goredis "github.com/go-redis/redis/v8"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/spf13/viper"
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

func NewCurrenciesService(ctx context.Context) (*CurrenciesService, error) {
	repo, err := db.NewCurrenciesRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.NewCurrenciesRepo: %w", err)
	}

	redisPool := &redis.Pool{
		MaxActive: 10,
		MaxIdle:   10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf(":%s", viper.GetString("redis.port")))
		},
	}
	enqueuer := work.NewEnqueuer("currency_rates", redisPool)
	pool := work.NewWorkerPool(ctx, 10, "currency_rates", redisPool)

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: "",
		DB:       0,
	})

	service := CurrenciesService{
		Repo:         repo,
		Enqueuer:     enqueuer,
		Pool:         pool,
		CurrencyList: domain.NewCurrencyList(),
		KV:           rdb,
	}

	service.Pool.JobWithOptions("update_currency_rate", work.JobOptions{MaxFails: 5}, service.UpdateRateJob)
	service.Pool.Start()

	return &service, nil
}

func (s *CurrenciesService) GetCode(code string) uint8 {
	return s.CurrencyList.AvailableCurrencies[code]
}
