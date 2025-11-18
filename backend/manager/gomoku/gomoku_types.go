package gomoku

import "boredgamz/manager"

type GomokuLobbyRequest struct {
	LobbyType string `json:"lobbyType"`
	Player  *manager.Player `json:"player"`
}

type GomokuClientRequest struct {
    Type string      `json:"type"`
    Data *GomokuGameState `json:"data"`
}

type GomokuServerResponse struct {
	Type string      `json:"type"`
  Data *GomokuGameState `json:"data"`
}

type GomokuMoveRequest struct {
    Row int `json:"row"`
    Col int `json:"col"`
    Player *manager.Player `json:"player"`
}
