package gomoku

import (
	"boredgamz/core/gomoku"
	"boredgamz/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func InsertGame(db *db.Database, gameState *gomoku.GomokuGameState) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	query := "INSERT INTO gomoku_games (id, player1_id, player2_id, game_state) VALUES ($1, $2, $3, $4)"
	gameStateBytes, err := json.Marshal(gameState)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	_, err = db.Pool.Exec(ctx, query, gameState.GameID, gameState.Players[0].PlayerID, gameState.Players[1].PlayerID, gameStateBytes)
	log.Println(err)
	if err != nil {
		return fmt.Errorf("failed to insert game: %w", err)
	}

	return nil
}

func GetGameByID(db *db.Database, gameID string) (*gomoku.GomokuGameState, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := "SELECT game_state FROM gomoku_games WHERE id=$1"

    var rawJSON []byte

    err := db.Pool.QueryRow(ctx, query, gameID).Scan(&rawJSON)
    if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
    if err != nil {
      return nil, fmt.Errorf("db query error: %w", err)
    }

    var gameState gomoku.GomokuGameState
    if err := json.Unmarshal(rawJSON, &gameState); err != nil {
        return nil, fmt.Errorf("json unmarshal error: %w", err)
    }

    return &gameState, nil
}


func GetGamesByPlayerID(db *db.Database, playerID string) ([]*gomoku.GomokuGameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT game_state FROM gomoku_games games WHERE games.player1_id=$1 OR games.player2_id=$1"
	rows, err := db.Pool.Query(ctx, query, playerID)
	
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("db query error: %w", err)
	}

	var games []*gomoku.GomokuGameState
	for rows.Next() {
		var gameState gomoku.GomokuGameState
		if err := rows.Scan(&gameState); err != nil {
			return nil, fmt.Errorf("failed to scan game state: %w", err)
		}
		games = append(games, &gameState)
	}


	return games, nil
}


