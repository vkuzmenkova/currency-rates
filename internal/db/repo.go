package db

import (
	"context"
	"fmt"

	"github.com/vkuzmenkova/currency-rates/configs"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type CurrenciesRepo struct {
	Conn DB
}

func NewCurrenciesRepo(ctx context.Context, config configs.Config) (*CurrenciesRepo, error) {
	conn, err := newPGConnection(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("PG connection failed.")
		return nil, fmt.Errorf("newPGConnection: %w", err)
	}
	log.Info().Msg("PG connection established")

	return &CurrenciesRepo{
		Conn: conn,
	}, nil
}

func newPGConnection(ctx context.Context, config configs.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		config.DB.Driver,
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.SSLMode,
	))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}
