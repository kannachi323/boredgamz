package arena

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGameRejectsUnknownKinds(t *testing.T) {
	_, err := NewGame("checkers")
	require.Error(t, err)
}

func TestTicTacToeWinsAndRejectsOccupiedMoves(t *testing.T) {
	game, err := NewGame(TicTacToe)
	require.NoError(t, err)

	for _, position := range []int{0, 1, 2} {
		actual, moveErr := game.ApplyMove(position, "one")
		require.NoError(t, moveErr)
		require.Equal(t, position, actual)
	}
	require.Equal(t, "one", game.Winner(2))
	_, err = game.ApplyMove(2, "two")
	require.EqualError(t, err, "position is already occupied")
}

func TestGomokuDetectsDiagonalFromEitherEnd(t *testing.T) {
	game, err := NewGame(Gomoku)
	require.NoError(t, err)
	for step := 0; step < 5; step++ {
		_, err = game.ApplyMove((3+step)*game.Columns+(8-step), "two")
		require.NoError(t, err)
	}
	require.Equal(t, "two", game.Winner(7*game.Columns+4))
	require.Len(t, game.gridMoves, 5)
	require.Equal(t, GridMove{Position: 7*game.Columns + 4, Mark: "two"}, game.gridMoves[4])
}

func TestConnectFourDropsToLowestAvailableRow(t *testing.T) {
	game, err := NewGame(ConnectFour)
	require.NoError(t, err)

	positions := []int{}
	for index := 0; index < game.Rows; index++ {
		position, moveErr := game.ApplyMove(3, "one")
		require.NoError(t, moveErr)
		positions = append(positions, position)
	}
	require.Equal(t, []int{38, 31, 24, 17, 10, 3}, positions)
	_, err = game.ApplyMove(3, "two")
	require.EqualError(t, err, "column is full")
}

func TestBoardDraw(t *testing.T) {
	game, err := NewGame(TicTacToe)
	require.NoError(t, err)
	moves := []struct {
		position int
		mark     string
	}{{0, "one"}, {1, "two"}, {2, "one"}, {3, "one"}, {4, "two"}, {5, "two"}, {6, "two"}, {7, "one"}, {8, "one"}}
	for _, move := range moves {
		_, err = game.ApplyMove(move.position, move.mark)
		require.NoError(t, err)
	}
	require.True(t, game.IsDraw())
	require.Empty(t, game.Winner(8))
}

func TestTicTacToeModesHaveAuthoritativeRules(t *testing.T) {
	fourSquare, err := NewGameWithMode(TicTacToe, "four-square")
	require.NoError(t, err)
	require.Equal(t, 4, fourSquare.Rows)
	require.Equal(t, 16, len(fourSquare.Board))
	require.Equal(t, 3, fourSquare.Connect)

	misere, err := NewGameWithMode(TicTacToe, "misere")
	require.NoError(t, err)
	require.True(t, misere.Misere)
}

func TestProGomokuRequiresExactlyFive(t *testing.T) {
	game, err := NewGameWithMode(Gomoku, "pro")
	require.NoError(t, err)
	for column := 0; column < 6; column++ {
		_, err = game.ApplyMove(7*game.Columns+column, "one")
		require.NoError(t, err)
	}
	require.Empty(t, game.Winner(7*game.Columns+5))
}

func TestChessUsesValidatedMovesAndTracksNotation(t *testing.T) {
	game, err := NewGameWithMode(Chess, "standard")
	require.NoError(t, err)
	require.Len(t, game.Board, 64)
	require.Equal(t, "br", game.Board[0])
	require.Equal(t, "wk", game.Board[60])
	require.EqualError(t, game.ApplyChessMove("e2", "e5", ""), "illegal chess move")
	require.NoError(t, game.ApplyChessMove("e2", "e4", ""))
	state := game.ChessState()
	require.Contains(t, state.FEN, " b KQkq ")
	require.Equal(t, "e4", state.Moves[0].SAN)
	require.Equal(t, "e2", state.LastFrom)
	require.Equal(t, "e4", state.LastTo)
}

func TestChessDetectsCheckmateAuthoritatively(t *testing.T) {
	game, err := NewGame(Chess)
	require.NoError(t, err)
	for _, move := range [][2]string{{"f2", "f3"}, {"e7", "e5"}, {"g2", "g4"}, {"d8", "h4"}} {
		require.NoError(t, game.ApplyChessMove(move[0], move[1], ""))
	}
	winner, finished := game.ChessResult()
	require.True(t, finished)
	require.Equal(t, "two", winner)
	require.Equal(t, "Checkmate", game.ChessState().Method)
}

func TestChessRematchAlwaysStartsWithWhite(t *testing.T) {
	white := newPlayer("white", "White")
	black := newPlayer("black", "Black")
	white.Mark, black.Mark = "one", "two"
	game, err := NewGame(Chess)
	require.NoError(t, err)
	room := newRoom("chess-room", game, white, black, RoomSettings{Mode: "standard", BestOf: 3, TimeControl: "untimed"}, nil)
	room.finishRoundLocked("one")
	require.NoError(t, room.ReadyRematch(white.ID))
	require.NoError(t, room.ReadyRematch(black.ID))
	require.Equal(t, 2, room.round)
	require.Equal(t, "one", room.turn)
	require.Contains(t, room.game.ChessState().FEN, " w ")
}

func TestChessDrawOfferCanBeDeclinedByMoving(t *testing.T) {
	white := newPlayer("white-draw", "White")
	black := newPlayer("black-draw", "Black")
	white.Mark, black.Mark = "one", "two"
	game, err := NewGame(Chess)
	require.NoError(t, err)
	room := newRoom("draw-room", game, white, black, RoomSettings{Mode: "standard", BestOf: 1, TimeControl: "untimed"}, nil)
	require.EqualError(t, room.OfferDraw(white.ID), "offer a draw after making your move")
	require.NoError(t, room.HandleChessMove(white.ID, "e2", "e4", ""))
	require.NoError(t, room.OfferDraw(white.ID))
	require.Equal(t, "one", room.Snapshot().DrawOfferedBy)
	require.NoError(t, room.HandleChessMove(black.ID, "e7", "e5", ""))
	require.Empty(t, room.Snapshot().DrawOfferedBy)
}

func TestChessDrawOfferCanBeAccepted(t *testing.T) {
	white := newPlayer("white-agree", "White")
	black := newPlayer("black-agree", "Black")
	white.Mark, black.Mark = "one", "two"
	game, err := NewGame(Chess)
	require.NoError(t, err)
	room := newRoom("agreed-draw", game, white, black, RoomSettings{Mode: "standard", BestOf: 1, TimeControl: "untimed"}, nil)
	require.NoError(t, room.HandleChessMove(white.ID, "e2", "e4", ""))
	require.NoError(t, room.OfferDraw(white.ID))
	require.EqualError(t, room.AcceptDraw(white.ID), "there is no opponent draw offer to accept")
	require.NoError(t, room.AcceptDraw(black.ID))
	snapshot := room.Snapshot()
	require.Equal(t, "round-finished", snapshot.Status)
	require.Empty(t, snapshot.Winner)
	require.Empty(t, snapshot.DrawOfferedBy)
}
