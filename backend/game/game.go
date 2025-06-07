package game

import "log"

type GameState struct {
	Board   *Board      `json:"board"`
	Players []*Player  `json:"players"`
	Turn    string     `json:"turn"`
	Status  *GameStatus    `json:"status"`
	LastMove *Move    `json:"lastMove"`
}

type GameStatus struct {
	Winner string `json:"winner"`
	Draw   bool   `json:"draw"`
	Status string `json:"status"`
}

func CreateGameState(size int, p1 *Player, p2 *Player) *GameState {

	newGameState := &GameState{
		Board:  NewEmptyBoard(size),
		Turn:   p1.PlayerID,
		Players: []*Player{p1, p2},
		Status: &GameStatus{
			Winner: "",
			Draw: false,
			Status: "online",
		},
		LastMove: nil,
	}

	return newGameState
}

func UpdateGameState(gameState *GameState, clientGameState *GameState) {
	//first, check the board to see if new move is valid

	if IsValidMove(clientGameState.Board, clientGameState.LastMove) {

		AddStoneToBoard(gameState.Board, clientGameState.LastMove, &Stone{Color: "white"})
		
		//next, check for a win only if move is valid
		if IsGomoku(gameState.Board.Stones, clientGameState.LastMove, clientGameState.LastMove.Color) {
			newStatus := &GameStatus{
				Winner: gameState.Turn,
				Draw:  false,
				Status: "offline",
			}
			gameState.Status = newStatus

			return
		}

		if IsDraw(gameState.Board) {
			newStatus := &GameStatus{
				Winner: "",
				Draw: true,
				Status: "offline",
			}
			gameState.Status = newStatus
			
			return
		}
		//if no win or draw, just update the turn
		if (gameState.Turn == "P1") {
			gameState.Turn = "P2"
		} else if (gameState.Turn == "P2") {
			gameState.Turn = "P1"
		}

		gameState.LastMove = clientGameState.LastMove

		log.Print(gameState.Board)
	} else {
		log.Printf("Invalid move by player %s at %v", gameState.Turn, clientGameState.LastMove)
		return
	}
}
