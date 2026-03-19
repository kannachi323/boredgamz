package connectfour

type Stone struct {
	Color string
}

type Board struct {
	Rows   int
	Cols   int
	Stones [][]*Stone
}

type Move struct {
	Row int `json:"row"`
	Col int `json:"col"`
	Color string `json:"color"`
}

// Create an empty Connect Four board
func NewEmptyBoard(rows, cols int) *Board {
	stones := make([][]*Stone, rows)
	for r := 0; r < rows; r++ {
		stones[r] = make([]*Stone, cols)
	}
	return &Board{
		Rows:   rows,
		Cols:   cols,
		Stones: stones,
	}
}

// Get the next available row for a given column
func GetNextAvailableRow(board *Board, col int) (int, bool) {
	for r := board.Rows - 1; r >= 0; r-- { // bottom-up
		if board.Stones[r][col] == nil {
			return r, true
		}
	}
	return -1, false // column full
}

// Place a stone on the board
func AddPieceToBoard(board *Board, move *Move, stone *Stone) {
	board.Stones[move.Row][move.Col] = stone
}

// Check if the board is full
func IsDraw(board *Board) bool {
	for c := 0; c < board.Cols; c++ {
		if board.Stones[0][c] == nil {
			return false
		}
	}
	return true
}

// Check for 4 in a row from last move
func IsConnectFour(stones [][]*Stone, moveCol int, playerID string) bool {
	row, ok := getLastMoveRow(stones, moveCol)
	if !ok {
		return false
	}
	color := stones[row][moveCol].Color
	return checkDirection(stones, row, moveCol, color, 0, 1) || // horizontal
		checkDirection(stones, row, moveCol, color, 1, 0) || // vertical
		checkDirection(stones, row, moveCol, color, 1, 1) || // diagonal /
		checkDirection(stones, row, moveCol, color, 1, -1)   // diagonal \
}

// Helper: find last stone row in a column
func getLastMoveRow(stones [][]*Stone, col int) (int, bool) {
	for r := 0; r < len(stones); r++ {
		if stones[r][col] != nil {
			return r, true
		}
	}
	return -1, false
}

// Check 4 in a row along a given direction
func checkDirection(stones [][]*Stone, row, col int, color string, dRow, dCol int) bool {
	count := 1

	// forward direction
	r, c := row+dRow, col+dCol
	for inBounds(stones, r, c) && stones[r][c] != nil && stones[r][c].Color == color {
		count++
		r += dRow
		c += dCol
	}

	// backward direction
	r, c = row-dRow, col-dCol
	for inBounds(stones, r, c) && stones[r][c] != nil && stones[r][c].Color == color {
		count++
		r -= dRow
		c -= dCol
	}

	return count >= 4
}

// Check bounds
func inBounds(stones [][]*Stone, r, c int) bool {
	return r >= 0 && r < len(stones) && c >= 0 && c < len(stones[0])
}
