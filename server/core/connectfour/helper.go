package connectfour

import "boredgamz/core"

func GetPlayerByColor(gs *ConnectFourGameState, color string) *core.Player {
	for _, player := range gs.Players {
		if player.Color == color {
			return player
		}
	}
	return nil
}

func GetOpponent(gs *ConnectFourGameState, playerID string) *core.Player {
	for _, player := range gs.Players {
		if player.PlayerID != playerID {
			return player
		}
	}
	return nil
}

func GetPlayerByID(gs *ConnectFourGameState, playerID string) *core.Player {
	for _, player := range gs.Players {
		if player.PlayerID == playerID {
			return player
		}
	}
	return nil
}
