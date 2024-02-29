package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=DB
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=KVDB
type KVDB interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type CurrenciesRepo struct {
	Conn *pgx.Conn
}

func NewCurrenciesRepo(ctx context.Context) (*CurrenciesRepo, error) {
	conn, err := newPGConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("newPGConnection: %w", err)
	}
	log.Println("PG connection established")

	return &CurrenciesRepo{
		Conn: conn,
	}, nil
}

func newPGConnection(ctx context.Context) (*pgx.Conn, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read db config: %w", err)
	}

	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("db.driver"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.name"),
		viper.GetString("db.ssl_mode"),
	))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}
