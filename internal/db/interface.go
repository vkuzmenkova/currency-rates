package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vkuzmenkova/currency-rates/models"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=Repo
type Repo interface {
	GetLastRate(ctx context.Context, base uint8, currencyCode uint8)
	GetRateByUUID(ctx context.Context, uuid uuid.UUID) (models.CurrencyRate, error)
	InsertRate(ctx context.Context, uuidUpdate string, base uint8, currencyCode uint8, value float64) error
	GetCurrencies(ctx context.Context) (map[string]uint8, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=DB
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=Row
type Row interface {
	Scan(dest ...any) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=KeyValueDB
type KeyValueDB interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=KVResult
type KVResult interface {
	Result() (string, error)
}
