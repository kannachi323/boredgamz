package manager

/*NOTES FOR MYSELF
this is a simple queue system that just puts last 2 players in a room
(given NumPlayers >= 2). Each board size will have a different lobby, and each
lobby will be in charge of matching players quickly
*/

import (
	"container/list"
	"fmt"
	"sync"

	"playgomoku/backend/game"
)

type Lobby struct {
	Queue  *list.List
	PlayerMap map[*game.Player]*list.Element
	NumPlayers	int
	MaxPlayers	int
	LobbyType	string
	RoomManager *RoomManager
}

type LobbyManager struct {
	lobbies map[string]*Lobby
	mu sync.RWMutex
}

func NewLobbyManager() *LobbyManager {
	return &LobbyManager{
		lobbies: make(map[string]*Lobby),
	}
}

func (lm *LobbyManager) GetLobby(lobbyType string) (*Lobby, bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lobby, ok := lm.lobbies[lobbyType]
	return lobby, ok
}

func (lm *LobbyManager) CreateLobby(maxPlayers int, lobbyType string) *Lobby {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lobby := &Lobby{
		Queue: list.New(),
		PlayerMap: make(map[*game.Player]*list.Element),
		NumPlayers: 0,
		MaxPlayers: maxPlayers,
		LobbyType: lobbyType,
		RoomManager: CreateRoomManager(),
	}
	lm.lobbies[lobbyType] = lobby

	return lobby
}

func (lm* LobbyManager) AddPlayerToQueue(lobby *Lobby, player *game.Player) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lobby.NumPlayers < lobby.MaxPlayers {
		slot := lobby.Queue.PushBack(player)
		lobby.PlayerMap[player] = slot
		lobby.NumPlayers++
	}
}

func (lm* LobbyManager) RemovePlayerFromQueue(lobby *Lobby, player *game.Player) {
	if elem, ok := lobby.PlayerMap[player]; ok {
		lobby.Queue.Remove(elem)
		delete(lobby.PlayerMap, player)
		lobby.NumPlayers--
	}
}

func (lm* LobbyManager) MatchPlayers(lobby *Lobby) (*Room, bool) {
	if lobby.NumPlayers >= 2 {
		e1 := lobby.Queue.Front()
		e2 := e1.Next()

		if e1 == nil || e2 == nil {
			return nil, false
		}

		player1 := e1.Value.(*game.Player)
		player2 := e2.Value.(*game.Player)

		// Remove both players from queue and player map
		lobby.Queue.Remove(e1)
		lobby.Queue.Remove(e2)
		delete(lobby.PlayerMap, player1)
		delete(lobby.PlayerMap, player2)
		lobby.NumPlayers -= 2

		fmt.Println("Matched players:", player1.PlayerID, player2.PlayerID)

		room := lobby.RoomManager.CreateNewRoom(player1, player2, lobby.LobbyType)

		player1.StartReader()
		player1.StartWriter()

		player2.StartReader()
		player2.StartWriter()

		return room, true
	}

	return nil, false
}
