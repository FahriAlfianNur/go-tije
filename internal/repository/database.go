package repository

import (
	"context"
	"fmt"

	"github.com/fahri/go-tije/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	
	return pool, nil
}