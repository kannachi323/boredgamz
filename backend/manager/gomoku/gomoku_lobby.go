package gomoku

import (
	"boredgamz/manager"
	"container/list"
	"fmt"
	"log"
	"sync"
)

type GomokuLobby struct {
	*manager.Lobby
	GomokuType string
	WhiteQueue *list.List
	BlackQueue *list.List
	PlayerMap  map[*manager.Player]*GomokuLobbySlot
	mu sync.RWMutex
}

type GomokuLobbySlot struct {
	Element *list.Element
	Queue *list.List
}

func NewGomokuLobby(maxPlayers int, gomokuType string) manager.LobbyController {
	gomokuLobby := &GomokuLobby{
		Lobby: &manager.Lobby{
			NumPlayers: 0,
			MaxPlayers: maxPlayers,
			RoomManager: nil,
		},
		GomokuType: gomokuType,
		WhiteQueue: list.New(),
		BlackQueue: list.New(),
		PlayerMap: make(map[*manager.Player]*GomokuLobbySlot),
	}
	return gomokuLobby
}

func (lobby *GomokuLobby) AddPlayer(player *manager.Player) {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()

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
	lobby.PlayerMap[player] = &GomokuLobbySlot{
		Element: elem,
		Queue: queue,
	}
	lobby.NumPlayers++
}

func (lobby *GomokuLobby) RemovePlayer(player *manager.Player) {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	slot, ok := lobby.PlayerMap[player]
	if !ok { return }

	slot.Queue.Remove(slot.Element)
	delete(lobby.PlayerMap, player)
	lobby.NumPlayers--
}

func (lobby *GomokuLobby) MatchPlayers() ([]*manager.Player, bool) {
	if lobby.WhiteQueue.Len() == 0 || lobby.BlackQueue.Len() == 0 { return nil, false }

	e1 := lobby.WhiteQueue.Front()
	e2 := lobby.BlackQueue.Front()
	playerWhite := e1.Value.(*manager.Player)
	playerBlack := e2.Value.(*manager.Player)

	lobby.RemovePlayer(playerWhite)
	lobby.RemovePlayer(playerBlack)

	if (playerWhite.PlayerID == playerBlack.PlayerID) { return nil, false }

	fmt.Println("Matched players:", playerWhite.PlayerID, playerBlack.PlayerID)

	return []*manager.Player{playerWhite, playerBlack}, true
}
