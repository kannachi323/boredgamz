package gomoku

import (
	"boredgamz/manager"
	"fmt"

	"github.com/google/uuid"
)

type GomokuGameStatus struct {
	Result string `json:"result"`
	Code string `json:"code"`
	Winner *manager.Player `json:"winner,omitempty"`
}

type GomokuGameState struct {
	GameID   string      `json:"gameID"`
	Board    *Board      `json:"board"`
	Players  []*manager.Player   `json:"players"`
	Status   *GomokuGameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
	Turn     string      `json:"turn"`
	Timeout  chan struct{} `json:"-"`
}


func NewGomokuGameState(gomokuType string, p1 *manager.Player, p2 *manager.Player) *GomokuGameState {
	var turn string
	if p1.Color == "black" {
		turn = p1.PlayerID
	} else {
		turn = p2.PlayerID
	}

	var size int
	switch gomokuType {
	case "19x19":
		size = 19
	case "15x15":
		size = 15
	case "9x9":
		size = 9
	default:
		size = 9
	}

	newGameState := &GomokuGameState{
		GameID:  uuid.New().String(),
		Board:   NewEmptyBoard(size),
		Players: []*manager.Player{p1, p2},
		Status: &GomokuGameStatus{
			Result: "",
			Code: "online",
			Winner: nil,
		},
		LastMove: nil,
		Turn: turn,
	}
	return newGameState
}

/*
Only UpdateGameState should be called by room handler
*/
func UpdateGameState(serverGameState *GomokuGameState, clientGameState *GomokuGameState) {
	
	err := UpdateLastMove(serverGameState, clientGameState.LastMove)
	if err != nil { return }

	if IsGomoku(serverGameState.Board.Stones, clientGameState.LastMove) {
		UpdateGameStatus(serverGameState, "win", serverGameState.Turn)
		return
	}

	if IsDraw(serverGameState.Board) {
		UpdateGameStatus(serverGameState, "draw", "")
		return
	}


	UpdatePlayerClocks(serverGameState);
	UpdatePlayerTurn(serverGameState);

	clientGameState = serverGameState; //VERY IMPORTANT: server game state is always source of truth
}

/*
Private handlers for updating game state
*/
func UpdatePlayerTurn(serverGameState *GomokuGameState) {
	switch serverGameState.Turn {
		case serverGameState.Players[0].PlayerID:
			serverGameState.Turn = serverGameState.Players[1].PlayerID
		case serverGameState.Players[1].PlayerID:
			serverGameState.Turn = serverGameState.Players[0].PlayerID
	}
}
func UpdateLastMove(serverGameState *GomokuGameState, move *Move) error {

	expectedColor := "black"
	if serverGameState.LastMove != nil {
		if serverGameState.LastMove.Color == "black" {
			expectedColor = "white"
		} else {
			expectedColor = "black"
		}
	}

	if move == nil || move.Color != expectedColor || !IsValidMove(serverGameState.Board, move) { return fmt.Errorf("invalid move") }
	
	AddStoneToBoard(serverGameState.Board, move, &Stone{Color: move.Color})
	serverGameState.LastMove = move

	return nil
}

func UpdateGameStatus(gs *GomokuGameState, statusType string, playerID string) {
	switch statusType {
	case "win":
		gs.Status = &GomokuGameStatus{
			Result: "win",
			Code:   "offline",
			Winner: GetPlayerByID(gs, playerID),
		}
	case "draw":
		gs.Status = &GomokuGameStatus{
			Result: "draw",
			Code:   "offline",
		}
	case "timeout":
		gs.Status = &GomokuGameStatus{
			Result: "win",
			Code:   "offline",
			Winner: GetOpponent(gs, GetPlayerByID(gs, playerID).Color),
		}
	}
}

func UpdatePlayerClocks(serverGameState *GomokuGameState) {
	currentPlayer := GetPlayerByColor(serverGameState, serverGameState.LastMove.Color)
	opponentPlayer := GetOpponent(serverGameState, serverGameState.LastMove.Color)

	currentPlayer.StopClock()
	opponentPlayer.StartClock()
}
