package game

import (
	"github.com/google/uuid"
)

type GameState struct {
	GameID   string      `json:"gameID"`
	Board    *Board      `json:"board"`
	Players  []*Player   `json:"players"`
	Turn     string      `json:"turn"`
	Status   *GameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
}

type GameStatus struct {
	Result string `json:"result"`
	Code string `json:"code"`
}

func CreateGameState(size int, p1 *Player, p2 *Player) *GameState {
	newGameState := &GameState{
		GameID:  uuid.New().String(),
		Board:   NewEmptyBoard(size),
		Turn:    p1.PlayerID,
		Players: []*Player{p1, p2},
		Status: &GameStatus{
			Result: "",
			Code: "online",
		},
		LastMove: nil,
	}

	return newGameState
}

/*
Only UpdateGameStateMove should be called by room handler
*/
func UpdateGameStateMove(serverGameState *GameState, clientGameState *GameState) {
	if IsValidMove(clientGameState.Board, clientGameState.LastMove) {

		if IsGomoku(serverGameState.Board.Stones, clientGameState.LastMove) {
			updateLastMove(serverGameState, clientGameState.LastMove)
			updateGameStatus(serverGameState, "win")
			return
		}

		if IsDraw(serverGameState.Board) {
			updateLastMove(serverGameState, clientGameState.LastMove)
			updateGameStatus(serverGameState, "draw")
			return
		}

		updateLastMove(serverGameState, clientGameState.LastMove)
		updatePlayerTurn(serverGameState)

		clientGameState = serverGameState; //VERY IMPORTANT: server game state is always source of truth
	}
}

/*
Private handlers for updating game state
*/
func updatePlayerTurn(serverGameState *GameState) {
	switch serverGameState.Turn {
		case "P1":
			serverGameState.Turn = "P2"
		case "P2":
			serverGameState.Turn = "P1"
	}
}

func updateLastMove(serverGameState *GameState, move *Move) {

	expectedColor := "black"
	if serverGameState.LastMove != nil {
		if serverGameState.LastMove.Color == "black" {
			expectedColor = "white"
		} else {
			expectedColor = "black"
		}
	}

	if move == nil || move.Color != expectedColor || !IsValidMove(serverGameState.Board, move) { return }
	
	AddStoneToBoard(serverGameState.Board, move, &Stone{Color: move.Color})
	serverGameState.LastMove = move
}

func updateGameStatus(serverGameState *GameState, statusType string) {
	switch statusType {
	case "win":
		newStatus := &GameStatus{
			Result: "win",
			Code: "offline",
		}
		serverGameState.Status = newStatus
	case "draw":
		newStatus := &GameStatus{
			Result: "draw",
			Code: "offline",
		}
		serverGameState.Status = newStatus
	}
}
