package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

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

    return nil
}

func (db *Database) GetUserByEmailPassword(email string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	var id string
	var hashedPassword string
	query := "SELECT id, password FROM users WHERE email = $1"
	
	err := db.Pool.QueryRow(ctx, query, email).Scan(&id, &hashedPassword)
	if err != nil {
		log.Println("error fetching user:", err)
		return "", fmt.Errorf("failed to get user by email: %w", err)
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("invalid password:", err)
		return "", fmt.Errorf("invalid password: %w", err)
	}

	return id, nil
}

func (db *Database) GetUserByID(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var email string
	query := "SELECT email FROM users WHERE id = $1"
	
	err := db.Pool.QueryRow(ctx, query, userID).Scan(&email)
	if errors.Is(err, pgx.ErrNoRows) {
    return "", fmt.Errorf("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user by ID: %w", err)
	}

	return email, nil
}
