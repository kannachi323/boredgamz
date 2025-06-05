package game

type GameState struct {
	Board   [][]*Stone `json:"board"`
	Size    int        `json:"size"`
	Players []*Player  `json:"players"`
	Turn    string     `json:"turn"`
	Status  string     `json:"status"`
}

func CreateGameState(size int, p1 *Player, p2 *Player) *GameState {

	newGameState := &GameState{
		Board:  NewEmptyBoard(size),
		Size:   size,
		Turn:   p1.ID,
		Status: "active",
	}

	return newGameState
}
