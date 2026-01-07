package gomoku

import (
	"encoding/json"
)

type GomokuLobbyData struct {
  Name      string `json:"name"`
  Mode      string  `json:"mode"`
  TimeControl string `json:"timeControl"`
  PlayerID string `json:"playerID"`
  PlayerColor string `json:"playerColor"`
  PlayerName string   `json:"playerName"`
}

type GomokuMoveData struct {
  Move  Move        `json:"move"`
}

type GomokuGameStateData struct {
  GameState *GomokuGameState `json:"gameState"`
}

type GomokuReconnectData struct {
  LobbyID string `json:"lobbyID"`
  PlayerID string `json:"playerID"`
}

type GomokuClientRequest struct {
  Type string      `json:"type"`
  Data json.RawMessage `json:"data"`
}

type GomokuServerResponse struct {
  Type string `json:"type"`
  Data json.RawMessage  `json:"data"`
}
