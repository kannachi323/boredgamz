package gomoku

import (
	"boredgamz/core"
	"encoding/json"
)

type GomokuLobbyData struct {
  LobbyType string       `json:"lobbyType"`
  Player    *core.Player `json:"player"`
}

type GomokuMoveData struct {
  Move  Move        `json:"move"`
}

type GomokuGameStateData struct {
  GameState *GomokuGameState `json:"gameState"`
}

type GomokuClientRequest struct {
  Type string      `json:"type"`
  Data json.RawMessage `json:"data"`
}

type GomokuServerResponse struct {
  Type string `json:"type"`
  Data json.RawMessage  `json:"data"`
}
