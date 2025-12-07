package connectfour

import (
	"boredgamz/core"
	"encoding/json"
)

type ConnectFourLobbyData struct {
	LobbyType string       `json:"lobbyType"`
	Player    *core.Player `json:"player"`
}

type ConnectFourMoveData struct {
	Column int `json:"column"` // Connect Four moves are column-based
}

type ConnectFourGameStateData struct {
	GameState *ConnectFourGameState `json:"gameState"`
}

type ConnectFourClientRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ConnectFourServerResponse struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
