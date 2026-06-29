package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
)

func NewPostgres(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), cfg.DSN())

	if err != nil {
		return nil, err
	}

	return pool, nil
}
