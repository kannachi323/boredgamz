package game

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
	R int `json:"r"`
	C int `json:"c"`
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
	log.Print("Adding stone to board at position:", move.R, move.C, "with color:", stone.Color)
	board.Stones[move.R][move.C] = stone
	board.NumStones++
}

func IsValidMove(board *Board, move *Move) bool {
	if move.R < 0 || move.R >= len(board.Stones) || move.C < 0 || move.C >= len(board.Stones[0]) {
		return false
	}

	return board.Stones[move.R][move.C].Color == ""
}

func IsDraw(board *Board) bool {
	return board.NumStones >= board.Size * board.Size
}

func IsGomoku(stones [][]*Stone, move *Move, color string) bool {
	directions := [][2]int{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
		{1, 1},
		{-1, -1},
		{1, -1},
		{-1, 1},
	}

	row := move.R
	col := move.C

	count := 0

	for _, d := range directions {
		dr, dc := d[0], d[1]
		for i := 0; i < 5; i++ {
			nr, nc := row + dr, col + dc
			if nr < 0 || nr >= len(stones) || nc < 0 || nc >= len(stones[0]) {
				break
			}
			if stones[nr][nc] != nil && stones[nr][nc].Color == move.Color {
				count++
			}
		}
		if count >= 5 {
			return true
		}
	}
	
	return false

}

