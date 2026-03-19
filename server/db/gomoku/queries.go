package gomoku

import (
	"boredgamz/core/gomoku/model"
	"boredgamz/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func InsertGame(db *db.Database, gameID string, player1ID string, player2ID string, gameState *model.GomokuGameStateRow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameStateBytes, err := json.Marshal(gameState)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	query := `INSERT INTO gomoku_games (id, player1_id, player2_id, game_state) VALUES ($1, $2, $3, $4)`
	_, err = db.Pool.Exec(ctx, query,
		gameID,
		player1ID,
		player2ID,
		gameStateBytes,
	)
	if err != nil {
		return fmt.Errorf("failed to insert game: %w", err)
	}

	return nil
}

// GetGameByID loads a single game by ID.
func GetGameByID(db *db.Database, gameID string) (*model.GomokuGameStateRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rawJSON []byte
	query := `SELECT game_state FROM gomoku_games WHERE id=$1`
	err := db.Pool.QueryRow(ctx, query, gameID).Scan(&rawJSON)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("db query error: %w", err)
	}

	var gameState model.GomokuGameStateRow
	if err := json.Unmarshal(rawJSON, &gameState); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &gameState, nil
}

// GetGamesByPlayerID loads all games for a player.
func GetGamesByPlayerID(db *db.Database, playerID string) ([]*model.GomokuGameStateRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT game_state FROM gomoku_games WHERE player1_id=$1 OR player2_id=$1`
	rows, err := db.Pool.Query(ctx, query, playerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("db query error: %w", err)
	}
	defer rows.Close()

	var games []*model.GomokuGameStateRow
	for rows.Next() {
		var rawJSON []byte
		if err := rows.Scan(&rawJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		var gameState model.GomokuGameStateRow
		if err := json.Unmarshal(rawJSON, &gameState); err != nil {
			return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
		}
		games = append(games, &gameState)
	}

	return games, nil
}
