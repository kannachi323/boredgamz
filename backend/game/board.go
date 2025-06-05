package game

type Stone struct {
    R     int       `json:"r"`
    C     int       `json:"c"`
    Color string    `json:"color"`
}

func NewEmptyBoard(boardSize int) [][]*Stone {
	board := make([][]*Stone, boardSize)
	for i := 0; i < boardSize; i++ {
		board[i] = make([]*Stone, boardSize)
	}
	return board;
}