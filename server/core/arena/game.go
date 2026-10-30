package arena

import (
	"errors"
	"fmt"
	"strings"

	chess "github.com/corentings/chess/v2"
)

type GameKind string

const (
	TicTacToe   GameKind = "tictactoe"
	ConnectFour GameKind = "connect-four"
	Gomoku      GameKind = "gomoku"
	Chess       GameKind = "chess"
)

type Game struct {
	Kind         GameKind `json:"kind"`
	Mode         string   `json:"mode"`
	Rows         int      `json:"rows"`
	Columns      int      `json:"columns"`
	Connect      int      `json:"connect"`
	ExactConnect bool     `json:"exactConnect"`
	Misere       bool     `json:"misere"`
	Board        []string `json:"board"`
	Moves        int      `json:"moves"`
	chessGame    *chess.Game
	chessMoves   []ChessMove
	gridMoves    []GridMove
	lastFrom     string
	lastTo       string
}

func NewGame(kind GameKind) (*Game, error) {
	return NewGameWithMode(kind, "")
}

func NormalizeMode(kind GameKind, mode string) string {
	switch kind {
	case TicTacToe:
		switch mode {
		case "four-square", "misere":
			return mode
		default:
			return "classic"
		}
	case ConnectFour:
		switch mode {
		case "quick", "connect-five":
			return mode
		default:
			return "classic"
		}
	case Gomoku:
		switch mode {
		case "compact", "pro":
			return mode
		default:
			return "freestyle"
		}
	case Chess:
		switch mode {
		case "rapid", "blitz":
			return mode
		default:
			return "standard"
		}
	default:
		return ""
	}
}

func NewGameWithMode(kind GameKind, mode string) (*Game, error) {
	mode = NormalizeMode(kind, mode)
	game := &Game{Kind: kind, Mode: mode}
	switch kind {
	case TicTacToe:
		switch mode {
		case "four-square":
			game.Rows, game.Columns, game.Connect = 4, 4, 3
		case "misere":
			game.Rows, game.Columns, game.Connect, game.Misere = 3, 3, 3, true
		default:
			game.Rows, game.Columns, game.Connect = 3, 3, 3
		}
	case ConnectFour:
		switch mode {
		case "quick":
			game.Rows, game.Columns, game.Connect = 5, 6, 4
		case "connect-five":
			game.Rows, game.Columns, game.Connect = 7, 8, 5
		default:
			game.Rows, game.Columns, game.Connect = 6, 7, 4
		}
	case Gomoku:
		if mode == "compact" {
			game.Rows, game.Columns, game.Connect = 13, 13, 5
		} else {
			game.Rows, game.Columns, game.Connect = 15, 15, 5
		}
		game.ExactConnect = mode == "pro"
	case Chess:
		game.Rows, game.Columns = 8, 8
		game.chessGame = chess.NewGame()
	default:
		return nil, fmt.Errorf("unsupported game %q", kind)
	}
	game.Board = make([]string, game.Rows*game.Columns)
	if kind == Chess {
		game.syncChessBoard()
	}
	return game, nil
}

func (g *Game) ApplyChessMove(from, to, promotion string) error {
	if g.Kind != Chess || g.chessGame == nil {
		return errors.New("this room is not a chess game")
	}
	from, to, promotion = strings.ToLower(strings.TrimSpace(from)), strings.ToLower(strings.TrimSpace(to)), strings.ToLower(strings.TrimSpace(promotion))
	if len(from) != 2 || len(to) != 2 {
		return errors.New("invalid chess coordinates")
	}
	wanted := from + to + promotion
	var selected *chess.Move
	var san string
	position := g.chessGame.Position()
	validMoves := g.chessGame.ValidMoves()
	for index := range validMoves {
		move := &validMoves[index]
		encoded := move.String()
		matches := encoded == wanted || (promotion == "" && strings.HasPrefix(encoded, from+to) && (len(encoded) == 4 || strings.HasSuffix(encoded, "q")))
		if matches {
			selected = move
			san = chess.AlgebraicNotation{}.Encode(position, move)
			break
		}
	}
	if selected == nil {
		return errors.New("illegal chess move")
	}
	if err := g.chessGame.Move(selected, nil); err != nil {
		return fmt.Errorf("illegal chess move: %w", err)
	}
	promo := ""
	if len(selected.String()) > 4 {
		promo = selected.String()[4:]
	}
	g.chessMoves = append(g.chessMoves, ChessMove{From: from, To: to, Promotion: promo, SAN: san})
	g.lastFrom, g.lastTo, g.Moves = from, to, g.Moves+1
	g.syncChessBoard()
	return nil
}

func (g *Game) ChessState() *ChessState {
	if g.Kind != Chess || g.chessGame == nil {
		return nil
	}
	return &ChessState{FEN: g.chessGame.Position().String(), Moves: append([]ChessMove{}, g.chessMoves...), LastFrom: g.lastFrom, LastTo: g.lastTo, Outcome: g.chessGame.Outcome().String(), Method: g.chessGame.Method().String()}
}

func (g *Game) ChessResult() (winner string, finished bool) {
	if g.Kind != Chess || g.chessGame == nil || g.chessGame.Outcome() == chess.NoOutcome {
		return "", false
	}
	switch g.chessGame.Outcome() {
	case chess.WhiteWon:
		return "one", true
	case chess.BlackWon:
		return "two", true
	default:
		return "", true
	}
}

func (g *Game) syncChessBoard() {
	g.Board = make([]string, 64)
	for square, piece := range g.chessGame.Position().Board().SquareMap() {
		index := (7-int(square.Rank()))*8 + int(square.File())
		color := "w"
		if piece.Color() == chess.Black {
			color = "b"
		}
		g.Board[index] = color + chessPieceCode(piece.Type())
	}
}

func chessPieceCode(piece chess.PieceType) string {
	switch piece {
	case chess.King:
		return "k"
	case chess.Queen:
		return "q"
	case chess.Rook:
		return "r"
	case chess.Bishop:
		return "b"
	case chess.Knight:
		return "n"
	default:
		return "p"
	}
}

func (g *Game) ApplyMove(position int, mark string) (int, error) {
	if mark != "one" && mark != "two" {
		return -1, errors.New("invalid player mark")
	}

	var target int
	switch g.Kind {
	case ConnectFour:
		if position < 0 || position >= g.Columns {
			return -1, errors.New("column is outside the board")
		}
		target = -1
		for row := g.Rows - 1; row >= 0; row-- {
			candidate := row*g.Columns + position
			if g.Board[candidate] == "" {
				target = candidate
				break
			}
		}
		if target == -1 {
			return -1, errors.New("column is full")
		}
	default:
		if position < 0 || position >= len(g.Board) {
			return -1, errors.New("position is outside the board")
		}
		target = position
		if g.Board[target] != "" {
			return -1, errors.New("position is already occupied")
		}
	}

	g.Board[target] = mark
	g.gridMoves = append(g.gridMoves, GridMove{Position: target, Mark: mark})
	g.Moves++
	return target, nil
}

func (g *Game) Winner(position int) string {
	if position < 0 || position >= len(g.Board) {
		return ""
	}
	mark := g.Board[position]
	if mark == "" {
		return ""
	}

	row, column := position/g.Columns, position%g.Columns
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, direction := range directions {
		count := 1
		count += g.countDirection(row, column, direction[0], direction[1], mark)
		count += g.countDirection(row, column, -direction[0], -direction[1], mark)
		if count >= g.Connect && (!g.ExactConnect || count == g.Connect) {
			return mark
		}
	}
	return ""
}

func (g *Game) IsDraw() bool {
	return g.Moves == len(g.Board)
}

func (g *Game) countDirection(row, column, rowStep, columnStep int, mark string) int {
	count := 0
	for {
		row += rowStep
		column += columnStep
		if row < 0 || row >= g.Rows || column < 0 || column >= g.Columns {
			return count
		}
		if g.Board[row*g.Columns+column] != mark {
			return count
		}
		count++
	}
}
