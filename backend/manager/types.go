package manager

import (
	"playgomoku/backend/game"
)

type ClientRequest struct {
    Type string      `json:"type"`
    Data game.GameState `json:"data"`
}

type ServerResponse struct {
	Type string      `json:"type"`
    Data *game.GameState `json:"data"`
}

type LobbyRequest struct {
    LobbyType string `json:"lobbyType"`
    Player  game.Player `json:"player"`
}



