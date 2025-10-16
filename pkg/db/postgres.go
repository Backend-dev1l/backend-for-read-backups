package pool

import (
	"context"
	"fmt"
	"log/slog"
	"test-http/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	DSN := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.DBName, cfg.Postgres.SSLMode, cfg.Postgres.Password)

	poolConfig, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		log.Error("Failed to parse DSN", slog.Any("error", err))
		return nil, err
	}

	poolConfig.MaxConns = cfg.PoolConfig.MaxConns
	poolConfig.MinConns = cfg.PoolConfig.MinConns
	poolConfig.MaxConnLifetime = cfg.PoolConfig.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.PoolConfig.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.PoolConfig.HealthCheckPeriod

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Error("Failed to create connection pool", slog.Any("error", err))
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Error("Error when pinging database", slog.Any("error", err))
		db.Close()
		return nil, err
	}
	return db, nil
}
