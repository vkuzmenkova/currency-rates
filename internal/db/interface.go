package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=DB
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
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
