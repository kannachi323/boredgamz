package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

func (db *Database) CreateUser(email string, password string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    query := "INSERT INTO users (name, email, password, is_admin) VALUES ($1, $2, $3, $4)"
    _, err = db.Pool.Exec(ctx, query, "Matthew", email, hashedPassword, false)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    log.Print("User created successfully\n")
    return nil
}
