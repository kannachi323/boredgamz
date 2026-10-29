package arena

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Table-driven rule checks for the grid games, covering every mode.
func TestGridModeDimensions(t *testing.T) {
	cases := []struct {
		kind                   GameKind
		mode                   string
		rows, columns, connect int
	}{
		{TicTacToe, "classic", 3, 3, 3},
		{TicTacToe, "four-square", 4, 4, 3},
		{TicTacToe, "misere", 3, 3, 3},
		{TicTacToe, "unknown", 3, 3, 3},
		{ConnectFour, "classic", 6, 7, 4},
		{ConnectFour, "quick", 5, 6, 4},
		{ConnectFour, "connect-five", 7, 8, 5},
		{Gomoku, "freestyle", 15, 15, 5},
		{Gomoku, "pro", 15, 15, 5},
		{Gomoku, "compact", 13, 13, 5},
	}
	for _, c := range cases {
		game, err := NewGameWithMode(c.kind, c.mode)
		require.NoError(t, err)
		require.Equal(t, [3]int{c.rows, c.columns, c.connect}, [3]int{game.Rows, game.Columns, game.Connect}, "%s/%s", c.kind, c.mode)
		require.Len(t, game.Board, c.rows*c.columns)
	}
}

func TestGridRejectsOutOfBoundsAndInvalidMarks(t *testing.T) {
	ttt, _ := NewGame(TicTacToe)
	_, err := ttt.ApplyMove(-1, "one")
	require.Error(t, err)
	_, err = ttt.ApplyMove(9, "one")
	require.Error(t, err)
	_, err = ttt.ApplyMove(0, "spectator")
	require.Error(t, err)

	c4, _ := NewGame(ConnectFour)
	_, err = c4.ApplyMove(7, "one")
	require.EqualError(t, err, "column is outside the board")
	_, err = c4.ApplyMove(-1, "one")
	require.Error(t, err)
}

func TestConnectFourWinsInEveryDirection(t *testing.T) {
	// Horizontal along the bottom row.
	game, _ := NewGame(ConnectFour)
	var last int
	for column := 0; column < 4; column++ {
		last, _ = game.ApplyMove(column, "one")
	}
	require.Equal(t, "one", game.Winner(last))

	// Vertical stack.
	game, _ = NewGame(ConnectFour)
	for i := 0; i < 4; i++ {
		last, _ = game.ApplyMove(2, "two")
	}
	require.Equal(t, "two", game.Winner(last))

	// Rising diagonal: build a staircase with filler discs underneath.
	game, _ = NewGame(ConnectFour)
	for column := 0; column < 4; column++ {
		for filler := 0; filler < column; filler++ {
			_, err := game.ApplyMove(column, "two")
			require.NoError(t, err)
		}
		last, _ = game.ApplyMove(column, "one")
	}
	require.Equal(t, "one", game.Winner(last))

	// Three in a row is not a win.
	game, _ = NewGame(ConnectFour)
	for column := 0; column < 3; column++ {
		last, _ = game.ApplyMove(column, "one")
	}
	require.Empty(t, game.Winner(last))
}

func TestConnectFiveNeedsFive(t *testing.T) {
	game, _ := NewGameWithMode(ConnectFour, "connect-five")
	var last int
	for column := 0; column < 4; column++ {
		last, _ = game.ApplyMove(column, "one")
	}
	require.Empty(t, game.Winner(last))
	last, _ = game.ApplyMove(4, "one")
	require.Equal(t, "one", game.Winner(last))
}

func TestFreestyleGomokuAcceptsOverline(t *testing.T) {
	game, _ := NewGameWithMode(Gomoku, "freestyle")
	var last int
	for column := 0; column < 6; column++ {
		last, _ = game.ApplyMove(7*game.Columns+column, "one")
	}
	require.Equal(t, "one", game.Winner(last))
}

func TestFourSquareConnectsThree(t *testing.T) {
	game, _ := NewGameWithMode(TicTacToe, "four-square")
	for _, position := range []int{5, 10} {
		_, err := game.ApplyMove(position, "two")
		require.NoError(t, err)
	}
	last, err := game.ApplyMove(15, "two")
	require.NoError(t, err)
	require.Equal(t, "two", game.Winner(last))
}

func TestMisereCompletingALineLoses(t *testing.T) {
	game, _ := NewGameWithMode(TicTacToe, "misere")
	one, two := newPlayer("one", "Ada"), newPlayer("two", "Bo")
	one.Mark, two.Mark = "one", "two"
	room := newRoom("misere-room", game, one, two, RoomSettings{Mode: "misere", BestOf: 1, TimeControl: "untimed"}, nil)
	for _, move := range []struct {
		player   string
		position int
	}{{"one", 0}, {"two", 3}, {"one", 1}, {"two", 4}, {"one", 2}} {
		require.NoError(t, room.HandleMove(move.player, move.position))
	}
	snapshot := room.Snapshot()
	require.Equal(t, "finished", snapshot.Status)
	require.Equal(t, "two", snapshot.Winner)
}
