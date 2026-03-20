package model

type Player struct {
	PlayerID   string `json:"playerID"`
	PlayerName string `json:"playerName"`
	Color      string `json:"color"`
}

type Move struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Color  string `json:"color"`
}

type ChatMessage struct {
	SenderID   string `json:"senderID"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
	SentAt     string `json:"sentAt"`
}

// ConnectFourGameStateRow represents a finished or persisted game
type ConnectFourGameStateRow struct {
	GameID   string         `json:"gameID"`
	Players  []*Player      `json:"players"`
	Moves    []*Move        `json:"moves"`
	Messages []*ChatMessage `json:"messages"`
	Result   string         `json:"result"`
	Winner   *Player        `json:"winner,omitempty"`
}
