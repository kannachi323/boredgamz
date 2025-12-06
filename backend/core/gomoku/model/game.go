package model


type Move struct {
    Row   int    `json:"row"`
    Col   int    `json:"col"`
    Color string `json:"color"`
}

type Player struct {
    PlayerID  string `json:"playerID"`
    PlayerName string `json:"playerName"`
    Color     string `json:"color"`
}

type GomokuGameStateRow struct {
    GameID string    `json:"gameID"`
    BoardSize int       `json:"boardSize"`
    Players []*Player `json:"players"`
    Moves  []*Move `json:"moves"`
    Result string    `json:"result"`
    Winner *Player    `json:"winner,omitempty"`
}


