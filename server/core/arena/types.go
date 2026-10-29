package arena

import "time"

const (
	disconnectGrace = 20 * time.Second
	roomRetention   = 5 * time.Minute
	maxChatHistory  = 60
	maxChatLength   = 240
)

type MatchKey struct {
	Kind        GameKind
	Mode        string
	BestOf      int
	TimeControl string
	InviteCode  string
}

type JoinOptions struct {
	Kind        GameKind
	PlayerID    string
	PlayerName  string
	Mode        string
	BestOf      int
	TimeControl string
	InviteCode  string
	Role        string
	RoomID      string
}

type RoomSettings struct {
	Mode        string `json:"mode"`
	BestOf      int    `json:"bestOf"`
	TimeControl string `json:"timeControl"`
	Visibility  string `json:"visibility"`
	InviteCode  string `json:"inviteCode,omitempty"`
}

type PlayerView struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Mark      string `json:"mark"`
	Connected bool   `json:"connected"`
}

type ChatMessage struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"senderId"`
	SenderName string    `json:"senderName"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Snapshot struct {
	RoomID        string           `json:"roomId"`
	Kind          GameKind         `json:"kind"`
	Status        string           `json:"status"`
	Rows          int              `json:"rows"`
	Columns       int              `json:"columns"`
	Board         []string         `json:"board"`
	Turn          string           `json:"turn,omitempty"`
	Winner        string           `json:"winner,omitempty"`
	SeriesWinner  string           `json:"seriesWinner,omitempty"`
	Players       [2]*PlayerView   `json:"players"`
	ViewerMark    string           `json:"viewerMark"`
	LastMove      int              `json:"lastMove"`
	Settings      RoomSettings     `json:"settings"`
	Round         int              `json:"round"`
	Score         map[string]int   `json:"score"`
	RematchReady  []string         `json:"rematchReady"`
	Chat          []ChatMessage    `json:"chat"`
	Spectators    int              `json:"spectators"`
	TurnDeadline  *time.Time       `json:"turnDeadline,omitempty"`
	Chess         *ChessState      `json:"chess,omitempty"`
	Clocks        map[string]int64 `json:"clocks,omitempty"`
	DrawOfferedBy string           `json:"drawOfferedBy,omitempty"`
	GridMoves     []GridMove       `json:"gridMoves"`
}

type ChessState struct {
	FEN      string      `json:"fen"`
	Moves    []ChessMove `json:"moves"`
	LastFrom string      `json:"lastFrom,omitempty"`
	LastTo   string      `json:"lastTo,omitempty"`
	Outcome  string      `json:"outcome,omitempty"`
	Method   string      `json:"method,omitempty"`
}

type ChessMove struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Promotion string `json:"promotion,omitempty"`
	SAN       string `json:"san"`
}

type GridMove struct {
	Position int    `json:"position"`
	Mark     string `json:"mark"`
}

type ClientMessage struct {
	Type      string `json:"type"`
	Position  int    `json:"position,omitempty"`
	Content   string `json:"content,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Promotion string `json:"promotion,omitempty"`
}

type ServerMessage struct {
	Type    string    `json:"type"`
	State   *Snapshot `json:"state,omitempty"`
	Message string    `json:"message,omitempty"`
}
