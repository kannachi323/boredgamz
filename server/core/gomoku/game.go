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
	OpeningRule string `json:"openingRule"`
	SwapRuleEnabled bool `json:"swapRuleEnabled"`
	FirstMoveCenterEnabled bool `json:"firstMoveCenterEnabled"`
	Players  []*core.Player   `json:"players"`
	PlayerClocks map[string]*core.PlayerClock `json:"-"`
	Status   *GomokuGameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
	Turn     string      `json:"turn"`
	Timeout  chan struct{} `json:"-"`
	Moves	[]*Move     `json:"moves"`
}

func NewGomokuGame(name string, p1 *core.Player, p2 *core.Player, openingRule string, swapRuleEnabled bool, firstMoveCenterEnabled bool) *GomokuGameState {
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
	case "13x13":
		size = 13
	case "9x9":
		size = 9
	default:
		size = 9
	}

	if openingRule == "" {
		openingRule = "freestyle"
	}

	newGameState := &GomokuGameState{
		GameID:  uuid.New().String(),
		Board:   NewEmptyBoard(size),
		OpeningRule: openingRule,
		SwapRuleEnabled: swapRuleEnabled,
		FirstMoveCenterEnabled: firstMoveCenterEnabled,
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

	if IsGomokuByRule(gs.Board.Stones, move, gs.OpeningRule) {
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

func HandleGomokuSwap(gs *GomokuGameState) error {
	if !gs.SwapRuleEnabled {
		return fmt.Errorf("swap rule is disabled")
	}

	if len(gs.Moves) != 1 {
		return fmt.Errorf("swap is only allowed after the first move")
	}

	current := GetPlayerByID(gs, gs.Turn)
	opponent := GetOpponent(gs, gs.Turn)
	if current == nil || opponent == nil {
		return fmt.Errorf("invalid players")
	}

	current.Color, opponent.Color = opponent.Color, current.Color
	UpdatePlayerTurn(gs)
	return nil
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

	if gs.FirstMoveCenterEnabled && len(gs.Moves) == 0 {
		center := gs.Board.Size / 2
		if move.Row != center || move.Col != center {
			return fmt.Errorf("first move must be at center")
		}
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
			Winner: GetOpponent(gs, playerID),
		}
	}
}
