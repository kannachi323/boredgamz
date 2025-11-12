package game

import (
	"fmt"

	"github.com/google/uuid"
)

type GameState struct {
	GameID   string      `json:"gameID"`
	Board    *Board      `json:"board"`
	Players  []*Player   `json:"players"`
	Status   *GameStatus `json:"status"`
	LastMove *Move       `json:"lastMove"`
	Turn     string      `json:"turn"`
	Timeout  chan struct{} `json:"-"`
}

type GameStatus struct {
	Result string `json:"result"`
	Code string `json:"code"`
	Winner string `json:"winner,omitempty"`
}

func CreateGameState(size int, p1 *Player, p2 *Player) *GameState {
	var turn string
	if p1.Color == "black" {
		turn = p1.PlayerID
	} else {
		turn = p2.PlayerID
	}

	newGameState := &GameState{
		GameID:  uuid.New().String(),
		Board:   NewEmptyBoard(size),
		Players: NewPlayers(p1, p2),
		Status: &GameStatus{
			Result: "",
			Code: "online",
			Winner: "",
		},
		LastMove: nil,
		Turn: turn,
	}
	return newGameState
}

/*
Only UpdateGameState should be called by room handler
*/
func UpdateGameState(serverGameState *GameState, clientGameState *GameState) {
	
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
func UpdatePlayerTurn(serverGameState *GameState) {
	switch serverGameState.Turn {
		case serverGameState.Players[0].PlayerID:
			serverGameState.Turn = serverGameState.Players[1].PlayerID
		case serverGameState.Players[1].PlayerID:
			serverGameState.Turn = serverGameState.Players[0].PlayerID
	}
}
func UpdateLastMove(serverGameState *GameState, move *Move) error {

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

func UpdateGameStatus(gs *GameState, statusType string, playerID string) {
	switch statusType {
	case "win":
		gs.Status = &GameStatus{
			Result: "win",
			Code:   "offline",
			Winner: playerID,
		}
	case "draw":
		gs.Status = &GameStatus{
			Result: "draw",
			Code:   "offline",
		}
	case "timeout":
		winner := GetOpponentPlayerByColor(gs, GetPlayerByID(gs, playerID).Color)
		gs.Status = &GameStatus{
			Result: "win",
			Code:   "offline",
			Winner: winner.PlayerID,
		}
	}
}

func UpdatePlayerClocks(serverGameState *GameState) {
	currentPlayer := GetPlayerByColor(serverGameState, serverGameState.LastMove.Color)
	opponentPlayer := GetOpponentPlayerByColor(serverGameState, serverGameState.LastMove.Color)

	currentPlayer.StopClock()
	opponentPlayer.StartClock()
}
