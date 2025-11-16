package manager

/*NOTES FOR MYSELF
this is a simple queue system that just puts last 2 players in a room
(given NumPlayers >= 2). Each board size will have a different lobby, and each
lobby will be in charge of matching players quickly
*/

import (
	"playgomoku/backend/game"
)

type LobbyController interface {
	AddPlayer(player *game.Player)
	MatchPlayers() ([]*game.Player, bool)
	RemovePlayer(player *game.Player)
}

type Lobby struct {
	NumPlayers  int
	MaxPlayers  int
	LobbyType   string
	RoomManager *RoomManager
}

type LobbyManager struct {
	Lobbies map[string]*Lobby
}