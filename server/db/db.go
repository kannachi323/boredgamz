package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

// Start connects using DATABASE_URL, applies pending migrations and keeps the
// pool. Pool size is bounded by DB_MAX_CONNS (default 10) so a single instance
// cannot exhaust a small managed Postgres plan.
func (db *Database) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w", err)
	}
	config.MaxConns = int32(envInt("DB_MAX_CONNS", 10))
	config.MinConns = 0
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := migrate(ctx, pool); err != nil {
		pool.Close()
		return err
	}
	db.Pool = pool
	return nil
}

func (db *Database) Ready(ctx context.Context) error {
	if db == nil || db.Pool == nil {
		return nil
	}
	return db.Pool.Ping(ctx)
}

func (db *Database) Stop() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func envInt(name string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
