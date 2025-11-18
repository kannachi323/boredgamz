package manager

/*NOTES FOR MYSELF
this is a simple queue system that just puts last 2 players in a room
(given NumPlayers >= 2). Each board size will have a different lobby, and each
lobby will be in charge of matching players quickly
*/

import (
	"sync"
)

type LobbyController interface {
	AddPlayer(player *Player)
	MatchPlayers() ([]*Player, bool)
	RemovePlayer(player *Player)
}

type Lobby struct {
	NumPlayers  int
	MaxPlayers  int
	RoomManager *RoomManager
}

//IMPORTANT: pass LobbyManager to server so all handlers have access
type LobbyManager struct {
	Lobbies map[string]LobbyController
	mu sync.RWMutex
}

func NewLobbyManager() *LobbyManager {
	return &LobbyManager{
		Lobbies: make(map[string]LobbyController),
	}
}

func (lm *LobbyManager) RegisterLobby(lobbyName string, lobby LobbyController) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.Lobbies[lobbyName] = lobby
}

func (lm *LobbyManager) DeactivateLobby(lobbyName string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	delete(lm.Lobbies, lobbyName)
}

func (lm *LobbyManager) GetLobby(lobbyName string) (LobbyController, bool) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	lobby, ok := lm.Lobbies[lobbyName]
	return lobby, ok
}