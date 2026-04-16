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
	PlayerClocks map[string]*core.PlayerClock `json:"-"`
	Status   *GomokuGameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
	Turn     string      `json:"turn"`
	Timeout  chan struct{} `json:"-"`
	Moves	[]*Move     `json:"moves"`
}

func NewGomokuGame(name string, p1 *core.Player, p2 *core.Player) *GomokuGameState {
	var turn string
	if p1.Color == "black" {
		turn = p1.PlayerID
	} else {
		turn = p2.PlayerID
	}

	var size int
	switch name {
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
		Moves: make([]*Move, 0),
	}

	return newGameState
}

/*HANDLERS*/
func HandleGomokuMove(gs *GomokuGameState, move *Move) {
    if err := UpdateLastMove(gs, move); err != nil { return }

    UpdateMoves(gs, move)

    if IsGomoku(gs.Board.Stones, move) {
        UpdateGameStatus(gs, "win", gs.Turn)
        return
    }

    if IsDraw(gs.Board) {
        UpdateGameStatus(gs, "draw", "")
        return
    }

    // Switch turn
    UpdatePlayerTurn(gs)

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

func UpdateLastMove(gs *GomokuGameState, move *Move) error {
	player := GetPlayerByColor(gs, move.Color)
    if gs.Turn != player.PlayerID {
        return fmt.Errorf("not your turn")
    }

    if !IsValidMove(gs.Board, move) {
        return fmt.Errorf("invalid move")
    }

    AddStoneToBoard(gs.Board, move, &Stone{Color: move.Color})
    gs.LastMove = move
    return nil
}

func UpdateMoves(serverGameState *GomokuGameState, move *Move) {
	serverGameState.Moves = append(serverGameState.Moves, move)
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
