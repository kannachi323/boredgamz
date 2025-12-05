package model


type Move struct {
    Row   int    `json:"row"`
    Col   int    `json:"col"`
    Color string `json:"color"`
}

type Player struct {
    PlayerID  string `json:"player_id"`
    PlayerName string `json:"player_name"`
    Color     string `json:"color"`
}

type GomokuGameStateRow struct {
    Moves  []*Move `json:"moves"`
    Turn   string    `json:"turn"`
    Result string    `json:"result"`
    Winner *Player    `json:"winner,omitempty"`
}


