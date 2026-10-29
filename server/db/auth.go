package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

// PasswordCost is the bcrypt work factor for new hashes. Existing hashes keep
// the cost they were created with; bcrypt reads it from the hash.
const PasswordCost = 12

// dummyHash lets failed lookups spend the same bcrypt time as real ones so
// response timing does not reveal which emails are registered.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("boredgamz-timing-equaliser"), PasswordCost)

var ErrEmailTaken = errors.New("email already registered")

type User struct {
	ID    string
	Name  string
	Email string
}

func (db *Database) CreateUser(name, email, password string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{Name: name, Email: email}
	query := "INSERT INTO users (name, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id"
	err = db.Pool.QueryRow(ctx, query, name, email, hashedPassword, false).Scan(&user.ID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, ErrEmailTaken
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (db *Database) GetUserByEmailPassword(email string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	var hashedPassword string
	query := "SELECT id, password FROM users WHERE email = $1"

	err := db.Pool.QueryRow(ctx, query, email).Scan(&id, &hashedPassword)
	if errors.Is(err, pgx.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		return "", errors.New("invalid credentials")
	}
	if err != nil {
		return "", fmt.Errorf("failed to query user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return id, nil
}

func (db *Database) GetUserByID(userID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &User{}
	query := "SELECT id, name, email FROM users WHERE id = $1"

	err := db.Pool.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
