package manager

import (
	"playgomoku/backend/game"
)

type ClientRequest struct {
    Type string      `json:"type"`
    Data *game.GameState `json:"data"`
}

type ServerResponse struct {
	Type string      `json:"type"`
    Data *game.GameState `json:"data"`
}

type LobbyRequest struct {
    LobbyType string `json:"lobbyType"`
    Player  *game.Player `json:"player"`
}

type MoveRequest struct {
    Row int `json:"row"`
    Col int `json:"col"`
    Player *game.Player `json:"player"`
}



