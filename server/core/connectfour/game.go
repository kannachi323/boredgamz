package connectfour

import (
	"boredgamz/core"
	"fmt"

	"github.com/google/uuid"
)

// Game status
type ConnectFourGameStatus struct {
	Result string       `json:"result"`
	Code   string       `json:"code"`
	Winner *core.Player `json:"winner,omitempty"`
}

// Game state
type ConnectFourGameState struct {
	GameID       string                   `json:"gameID"`
	Board        *Board                   `json:"board"`
	Players      []*core.Player           `json:"players"`
	PlayerClocks map[string]*core.PlayerClock `json:"-"`
	Status       *ConnectFourGameStatus   `json:"status"`
	LastMove     *Move                    `json:"lastMove"`
	Turn         string                   `json:"turn"`
	Timeout      chan struct{}            `json:"-"`
	Moves        []*Move                  `json:"moves"`
}

// New game
func NewConnectFourGame(gameType string, p1 *core.Player, p2 *core.Player) *ConnectFourGameState {
	var turn string
	if p1.Color == "red" { // Red goes first in Connect Four
		turn = p1.PlayerID
	} else {
		turn = p2.PlayerID
	}

	newGameState := &ConnectFourGameState{
		GameID:  uuid.New().String(),
		Board:   NewEmptyBoard(6, 7), // Connect Four standard 6 rows x 7 cols
		Players: []*core.Player{p1, p2},
		Status: &ConnectFourGameStatus{
			Result: "",
			Code:   "online",
			Winner: nil,
		},
		LastMove: nil,
		Turn:     turn,
		Moves:    make([]*Move, 0),
	}

	return newGameState
}

/* HANDLE MOVES */
func HandleConnectFourMove(gs *ConnectFourGameState, col int) {
	if err := UpdateLastMove(gs, col); err != nil {
		return
	}

	UpdateMoves(gs, col)

	if IsConnectFour(gs.Board.Stones, col, gs.Turn) {
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

/* PRIVATE UPDATERS */
func UpdatePlayerTurn(gs *ConnectFourGameState) {
	switch gs.Turn {
	case gs.Players[0].PlayerID:
		gs.Turn = gs.Players[1].PlayerID
	case gs.Players[1].PlayerID:
		gs.Turn = gs.Players[0].PlayerID
	}
}

func UpdateLastMove(gs *ConnectFourGameState, col int) error {
	player := GetPlayerByID(gs, gs.Turn)
	row, ok := GetNextAvailableRow(gs.Board, col)
	if !ok {
		return fmt.Errorf("invalid move: col full")
	}

	move := &Move{
		Row:   row,
		Col: col,
		Color: player.Color,
	}

	AddPieceToBoard(gs.Board, move, &Stone{Color: player.Color})
	gs.LastMove = move
	return nil
}

func UpdateMoves(gs *ConnectFourGameState, col int) {
	gs.Moves = append(gs.Moves, gs.LastMove)
}

func UpdateGameStatus(gs *ConnectFourGameState, statusType string, playerID string) {
	switch statusType {
	case "win":
		gs.Status = &ConnectFourGameStatus{
			Result: "win",
			Code:   "offline",
			Winner: GetPlayerByID(gs, playerID),
		}
	case "draw":
		gs.Status = &ConnectFourGameStatus{
			Result: "draw",
			Code:   "offline",
		}
	case "timeout":
		gs.Status = &ConnectFourGameStatus{
			Result: "win",
			Code:   "offline",
			Winner: GetOpponent(gs, GetPlayerByID(gs, playerID).Color),
		}
	}
}
