package gomoku

import "log"

type Board struct {
	Stones [][]*Stone `json:"stones"`
	Size    int   `json:"size"`
	NumStones int `json:"numStones"`
}

type Stone struct {
    Color string    `json:"color"`
}

type Move struct {
	Row int `json:"row"`
	Col int `json:"col"`
	Color string `json:"color"`
}

func NewEmptyBoard(size int) *Board {
	stones := make([][]*Stone, size)
	for i := 0; i < size; i++ {
		stones[i] = make([]*Stone, size)
		for j := 0; j < size; j++ {
			stones[i][j] = &Stone{Color: ""}
		}
	}

	emptyBoard := &Board{
		Stones: stones,
		Size: size,
		NumStones: 0,
	}

	
	return emptyBoard;
}

func AddStoneToBoard(board *Board, move *Move, stone *Stone) {
	log.Print("Adding stone to board at position:", move.Row, move.Col, "with color:", stone.Color)
	board.Stones[move.Row][move.Col] = stone
	board.NumStones++
}

func IsValidMove(board *Board, move *Move) bool {
	if move.Row < 0 || move.Row >= len(board.Stones) || move.Col < 0 || move.Col >= len(board.Stones[0]) {
		return false
	}

	return board.Stones[move.Row][move.Col].Color == "" //empty slot
}

func IsDraw(board *Board) bool {
	return board.NumStones >= board.Size * board.Size
}

func IsGomoku(stones [][]*Stone, move *Move) bool {
	directions := [][2]int{
		{0, 1},
		{1, 0}, 
		{1, 1},
		{1, -1},
	}

	row := move.Row
	col := move.Col
	var count int

	for _, d := range directions {
		dr, dc := d[0], d[1]
		count = 1
		for i := 1; i < 5; i++ {
			nr, nc := row + dr * i, col + dc * i
			if (nr < 0 || nr >= len(stones) || nc < 0 || nc >= len(stones[0])) { break }
			if (stones[nr][nc] == nil || stones[nr][nc].Color != move.Color) { break }
			count++
		}
		for i := 1; i < 5; i++ {
			nr, nc := row - dr * i, col - dc * i
			if (nr < 0 || nr >= len(stones) || nc < 0 || nc >= len(stones[0])) { break }
			if (stones[nr][nc] == nil || stones[nr][nc].Color != move.Color) { break }
			count++	
		}
		if count >= 5 {
			return true
		}
	}
	return false
}

