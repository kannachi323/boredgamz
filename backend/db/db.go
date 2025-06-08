package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) Start() error {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")

	fmt.Print(dsn)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Optional: Ping to ensure DB is reachable
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	db.Pool = pool

	return nil
}

func (db *Database) Stop() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}