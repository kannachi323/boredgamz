package db

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/schema.sql
var schemaSQL string

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dsn := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Optional: Ping to ensure DB is reachable
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		pool.Close()
		return fmt.Errorf("unable to apply database schema: %w", err)
	}

	db.Pool = pool

	return nil
}

func (db *Database) Stop() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
