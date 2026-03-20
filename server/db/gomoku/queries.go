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
	"github.com/jackc/pgx/v5/pgconn"
)

const gameType = "gomoku"

func InsertGame(db *db.Database, gameID string, player1ID string, player2ID string, gameState *model.GomokuGameStateRow) error {
	return SaveGameWithMessages(db, gameID, player1ID, player2ID, gameState, gameState.Messages)
}

func SaveGameWithMessages(db *db.Database, gameID string, player1ID string, player2ID string, gameState *model.GomokuGameStateRow, messages []*model.ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameStateCopy := *gameState
	gameStateCopy.Messages = nil

	gameStateBytes, err := json.Marshal(&gameStateCopy)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO games (id, game_type, player1_id, player2_id, game_state) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, query,
		gameID,
		gameType,
		player1ID,
		player2ID,
		gameStateBytes,
	)
	if err != nil {
		return fmt.Errorf("failed to insert game: %w", err)
	}

	if err := insertMessages(ctx, tx, gameID, messages); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetGameByID loads a single game by ID.
func GetGameByID(db *db.Database, gameID string) (*model.GomokuGameStateRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rawJSON []byte
	query := `SELECT game_state FROM games WHERE id=$1 AND game_type=$2`
	err := db.Pool.QueryRow(ctx, query, gameID, gameType).Scan(&rawJSON)
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

	messages, err := getMessagesByGameID(ctx, db.Pool, gameID)
	if err != nil {
		return nil, err
	}
	if len(messages) > 0 {
		gameState.Messages = messages
	}

	return &gameState, nil
}

// GetGamesByPlayerID loads all games for a player.
func GetGamesByPlayerID(db *db.Database, playerID string) ([]*model.GomokuGameStateRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, game_state FROM games WHERE game_type=$1 AND (player1_id=$2 OR player2_id=$2)`
	rows, err := db.Pool.Query(ctx, query, gameType, playerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("db query error: %w", err)
	}
	defer rows.Close()

	var games []*model.GomokuGameStateRow
	for rows.Next() {
		var gameID string
		var rawJSON []byte
		if err := rows.Scan(&gameID, &rawJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		var gameState model.GomokuGameStateRow
		if err := json.Unmarshal(rawJSON, &gameState); err != nil {
			return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
		}

		messages, err := getMessagesByGameID(ctx, db.Pool, gameID)
		if err != nil {
			return nil, err
		}
		if len(messages) > 0 {
			gameState.Messages = messages
		}

		games = append(games, &gameState)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gomoku rows error: %w", err)
	}

	return games, nil
}

type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type dbQuerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func insertMessages(ctx context.Context, exec dbExecutor, gameID string, messages []*model.ChatMessage) error {
	if len(messages) == 0 {
		return nil
	}

	query := `INSERT INTO messages (game_id, game_type, sender_id, sender_name, content, sent_at) VALUES ($1, $2, $3, $4, $5, $6)`
	for _, message := range messages {
		if message == nil {
			continue
		}

		sentAt, err := time.Parse(time.RFC3339, message.SentAt)
		if err != nil {
			sentAt = time.Now().UTC()
		}

		if _, err := exec.Exec(ctx, query, gameID, gameType, message.SenderID, message.SenderName, message.Content, sentAt); err != nil {
			return fmt.Errorf("failed to insert gomoku message: %w", err)
		}
	}

	return nil
}

func getMessagesByGameID(ctx context.Context, queryer dbQuerier, gameID string) ([]*model.ChatMessage, error) {
	query := `SELECT sender_id, sender_name, content, sent_at FROM messages WHERE game_id=$1 AND game_type=$2 ORDER BY sent_at ASC, id ASC`
	rows, err := queryer.Query(ctx, query, gameID, gameType)
	if err != nil {
		return nil, fmt.Errorf("failed to query gomoku messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		var (
			senderID   string
			senderName string
			content    string
			sentAt     time.Time
		)

		if err := rows.Scan(&senderID, &senderName, &content, &sentAt); err != nil {
			return nil, fmt.Errorf("failed to scan gomoku message: %w", err)
		}

		messages = append(messages, &model.ChatMessage{
			SenderID:   senderID,
			SenderName: senderName,
			Content:    content,
			SentAt:     sentAt.UTC().Format(time.RFC3339),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gomoku message rows error: %w", err)
	}

	return messages, nil
}
