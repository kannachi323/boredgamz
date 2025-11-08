package game

import "github.com/google/uuid"

type GameState struct {
	GameID   string      `json:"gameID"`
	Board    *Board      `json:"board"`
	Players  []*Player   `json:"players"`
	Turn     string      `json:"turn"`
	Status   *GameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
}

type GameStatus struct {
	Winner string `json:"winner"`
	Draw   bool   `json:"draw"`
	Status string `json:"status"`
}

func CreateGameState(size int, p1 *Player, p2 *Player) *GameState {

	newGameState := &GameState{
		GameID:  uuid.New().String(),
		Board:   NewEmptyBoard(size),
		Turn:    p1.PlayerID,
		Players: []*Player{p1, p2},
		Status: &GameStatus{
			Winner: "",
			Draw:   false,
			Status: "online",
		},
		LastMove: nil,
	}

	return newGameState
}

func UpdateGameState(serverGameState *GameState, clientGameState *GameState) {
	if IsValidMove(clientGameState.Board, clientGameState.LastMove) {

		AddStoneToBoard(serverGameState.Board, clientGameState.LastMove, &Stone{Color: "white"})

		if IsGomoku(serverGameState.Board.Stones, clientGameState.LastMove, clientGameState.LastMove.Color) {
			newStatus := &GameStatus{
				Winner: serverGameState.Turn,
				Draw:   false,
				Status: "offline",
			}
			serverGameState.Status = newStatus

			return
		}

		if IsDraw(serverGameState.Board) {
			newStatus := &GameStatus{
				Winner: "",
				Draw:   true,
				Status: "offline",
			}
			serverGameState.Status = newStatus

			return
		}

		switch serverGameState.Turn {
		case "P1":
			serverGameState.Turn = "P2"
		case "P2":
			serverGameState.Turn = "P1"
		}
		serverGameState.LastMove = clientGameState.LastMove
	}
}
