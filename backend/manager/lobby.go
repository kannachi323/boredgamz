package manager

/*NOTES FOR MYSELF
this is a simple queue system that just puts last 2 players in a room
(given NumPlayers >= 2). Each board size will have a different lobby, and each
lobby will be in charge of matching players quickly
*/

import (
	"container/list"
	"fmt"
	"log"
	"sync"

	"playgomoku/backend/game"
)

type Lobby struct {
	WhiteQueue  *list.List
	BlackQueue 	*list.List
	PlayerMap map[*game.Player]*LobbySlot
	NumPlayers	int
	MaxPlayers	int
	LobbyType	string
	RoomManager *RoomManager
}

type LobbySlot struct {
	Element *list.Element
	Queue *list.List
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
		WhiteQueue: list.New(),
		BlackQueue: list.New(),
		PlayerMap: make(map[*game.Player]*LobbySlot),
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

	if lobby.NumPlayers >= lobby.MaxPlayers { return }
	if _, exists := lobby.PlayerMap[player]; exists { return }

	var elem *list.Element
	var queue *list.List
	switch player.Color {
	case "white":
		log.Println("Adding player to white queue:", player.PlayerID)
		elem = lobby.WhiteQueue.PushBack(player)
		queue = lobby.WhiteQueue
	case "black":
		log.Println("Adding player to black queue:", player.PlayerID)
		elem = lobby.BlackQueue.PushBack(player)
		queue = lobby.BlackQueue
	default:
		return
	}
	lobby.PlayerMap[player] = &LobbySlot{
		Element: elem,
		Queue: queue,
	}
	lobby.NumPlayers++
}

func (lm* LobbyManager) RemovePlayerFromQueue(lobby *Lobby, player *game.Player) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	slot, ok := lobby.PlayerMap[player]
	if !ok { return }

	slot.Queue.Remove(slot.Element)
	delete(lobby.PlayerMap, player)
	lobby.NumPlayers--
}

func (lm* LobbyManager) MatchPlayers(lobby *Lobby) ([]*game.Player, bool) {
	if lobby.WhiteQueue.Len() == 0 || lobby.BlackQueue.Len() == 0 { return nil, false }

	e1 := lobby.WhiteQueue.Front()
	e2 := lobby.BlackQueue.Front()
	playerWhite := e1.Value.(*game.Player)
	playerBlack := e2.Value.(*game.Player)

	lm.RemovePlayerFromQueue(lobby, playerWhite)
	lm.RemovePlayerFromQueue(lobby, playerBlack)

	if (playerWhite.PlayerID == playerBlack.PlayerID) { return nil, false }

	fmt.Println("Matched players:", playerWhite.PlayerID, playerBlack.PlayerID)

	return []*game.Player{playerWhite, playerBlack}, true
}
