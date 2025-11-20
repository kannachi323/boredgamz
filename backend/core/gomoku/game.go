package gomoku

import (
	"boredgamz/core"
	"fmt"

	"github.com/google/uuid"
)

type GomokuGameStatus struct {
	Result string `json:"result"`
	Code string `json:"code"`
	Winner *core.Player `json:"winner,omitempty"`
}

type GomokuGameState struct {
	GameID   string      `json:"gameID"`
	Board    *Board      `json:"board"`
	Players  []*core.Player   `json:"players"`
	Status   *GomokuGameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
	Turn     string      `json:"turn"`
	Timeout  chan struct{} `json:"-"`
}


func NewGomokuGame(gomokuType string, p1 *core.Player, p2 *core.Player) *GomokuGameState {
	var turn string
	if p1.Color == "black" {
		turn = p1.PlayerID
		p1.StartClock()
	} else {
		turn = p2.PlayerID
		p2.StartClock()
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
		Players: []*core.Player{p1, p2},
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

/*HANDLERS*/
func HandleGomokuMove(serverGameState *GomokuGameState, row int, col int, color string) {
	move := &Move{
        Row: row,
        Col: col,
        Color: color,
    }

    err := UpdateLastMove(serverGameState, move)
	if err != nil { return }
    

    if IsGomoku(serverGameState.Board.Stones, move) {
        UpdateGameStatus(serverGameState, "win", serverGameState.Turn)
    } else if IsDraw(serverGameState.Board) {
        UpdateGameStatus(serverGameState, "draw", "")
    } else {
        UpdatePlayerClocks(serverGameState)
        UpdatePlayerTurn(serverGameState)
    }
}

/*
PRIVATE gamestate updaters
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
	currentPlayer := GetPlayerByID(serverGameState, serverGameState.Turn)
	opponentPlayer := GetOpponent(serverGameState, serverGameState.Turn)

	currentPlayer.StopClock()
	opponentPlayer.StartClock()
}
