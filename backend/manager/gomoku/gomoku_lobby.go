package gomoku

import (
	"container/list"
	"fmt"
	"log"
	"playgomoku/backend/game"
	"playgomoku/backend/manager"
	"sync"
)

type GomokuLobby struct {
	manager.LobbyController
	
	manager.Lobby
	manager.LobbyManager

}

type GomokuLobbySlot struct {
	Element *list.Element
	Queue *list.List
}

type GomokuLobbyManager struct {
	lobbies map[string]*GomokuLobby
	mu sync.RWMutex
}


func (lm *GomokuLobbyManager) GetLobby(lobbyType string) (*GomokuLobby, bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lobby, ok := lm.lobbies[lobbyType]
	return lobby, ok
}

func (lm *GomokuLobbyManager) CreateLobby(maxPlayers int, lobbyType string) *GomokuLobby {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lobby := &Lobby{
		WhiteQueue: list.New(),
		BlackQueue: list.New(),
		PlayerMap: make(map[*game.Player]*GomokuLobbySlot),
		NumPlayers: 0,
		MaxPlayers: maxPlayers,
		LobbyType: lobbyType,
		RoomManager: CreateRoomManager(),
	}
	lm.lobbies[lobbyType] = lobby

	return lobby
}

func (lm* GomokuLobbyManager) AddPlayerToQueue(lobby *GomokuLobby, player *game.Player) {
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

func (lm* GomokuLobbyManager) RemovePlayerFromQueue(lobby *GomokuLobby, player *game.Player) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	slot, ok := lobby.PlayerMap[player]
	if !ok { return }

	slot.Queue.Remove(slot.Element)
	delete(lobby.PlayerMap, player)
	lobby.NumPlayers--
}

func (lm* GomokuLobbyManager) MatchPlayers(lobby *GomokuLobby) ([]*game.Player, bool) {
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

