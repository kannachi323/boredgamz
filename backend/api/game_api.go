package api

import (
	"encoding/json"
	"net/http"

	"playgomoku/backend/game"
)

type NewGameStateRequest struct {
    Size    int     `json:"size"`
    P1      game.Player  `json:"p1"`
    P2      game.Player  `json:"p2"`
}

func NewGameState(w http.ResponseWriter, r *http.Request) {
	var reqBody NewGameStateRequest
	json.NewDecoder(r.Body).Decode(&reqBody)

	size := reqBody.Size
	p1 := reqBody.P1
	p2 := reqBody.P2


	newGameState := game.CreateGameState(size, &p1, &p2)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGameState)
}

