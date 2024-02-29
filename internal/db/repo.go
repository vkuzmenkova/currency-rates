package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"log"
)

type CurrenciesRepo struct {
	Conn DB
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
