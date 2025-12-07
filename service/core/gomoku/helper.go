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

func GetLobbyIdentifier(boardSize int, gameMode string, timeControl string) string {
    return fmt.Sprintf("%s-%dx%d-%s", gameMode, boardSize, boardSize, timeControl)
}

func GetSimpleLobbyIdentifier(boardSize int) string {
	return fmt.Sprintf("%dx%d", boardSize, boardSize)
}