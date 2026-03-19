package gomoku

import (
	"boredgamz/core"
	"fmt"
)

func GetPlayerByColor(gs *GomokuGameState, color string) *core.Player {
	for _, player := range gs.Players {
		if player.Color == color {
			return player
		}
	}
	return nil
}

func GetOpponent(gs *GomokuGameState, playerID string) *core.Player {
	for _, player := range gs.Players {
		if player.PlayerID != playerID {
			return player
		}
	}
	return nil
}

func GetPlayerByID(gs *GomokuGameState, playerID string) *core.Player {
	for _, player := range gs.Players {
		if player.PlayerID == playerID {
			return player
		}
	}
	return nil
}

func GetGomokuLobbyID(name string, mode string) string {
  return fmt.Sprintf("gomoku-%s-%s", mode, name)
}
