package connectfour

import "encoding/json"

type ConnectFourLobbyData struct {
	TimeControl string `json:"timeControl"`
	PlayerID    string `json:"playerID"`
	PlayerColor string `json:"playerColor"`
	PlayerName  string `json:"playerName"`
}

type ConnectFourMoveData struct {
	Column int `json:"column"` // Connect Four moves are column-based
}

type ConnectFourChatData struct {
	Content string `json:"content"`
}

type ConnectFourReconnectData struct {
	PlayerID string `json:"playerID"`
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
