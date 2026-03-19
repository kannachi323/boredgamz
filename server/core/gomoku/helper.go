package gomoku

import (
	"boredgamz/core"
	"fmt"
	"strings"
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

func GetGomokuLobbyID(name string, mode string, openingRule string) string {
	if openingRule == "" {
		openingRule = "freestyle"
	}
	return fmt.Sprintf("gomoku-%s-%s-%s", mode, name, strings.ToLower(openingRule))
}

func ParseOpeningRuleFromLobbyID(lobbyID string) string {
	parts := strings.Split(lobbyID, "-")
	if len(parts) < 4 {
		return "freestyle"
	}

	rule := strings.ToLower(parts[len(parts)-1])
	switch rule {
	case "freestyle", "standard", "renju":
		return rule
	default:
		return "freestyle"
	}
}

func NormalizeLobbyID(lobbyID string) string {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(lobbyID)), "-")
	if len(parts) < 3 {
		return lobbyID
	}

	if len(parts) == 3 && parts[0] == "gomoku" {
		return lobbyID + "-freestyle"
	}

	if len(parts) >= 4 && parts[0] == "gomoku" {
		rule := parts[len(parts)-1]
		switch rule {
		case "freestyle", "standard", "renju":
			return lobbyID
		default:
			return strings.Join(parts[:3], "-") + "-freestyle"
		}
	}

	return lobbyID
}
