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

// ConnectFourGameStateRow represents a finished or persisted game
type ConnectFourGameStateRow struct {
	GameID  string    `json:"gameID"`
	Players []*Player `json:"players"`
	Moves   []*Move   `json:"moves"`
	Result  string    `json:"result"`        // "win" or "draw"
	Winner  *Player   `json:"winner,omitempty"` // nil if draw
}