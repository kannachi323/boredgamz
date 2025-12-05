package gomoku

import (
	"boredgamz/core"
	"boredgamz/db"
	"container/list"
	"log"
	"sync"
	"time"
)

type GomokuLobby struct {
	*core.Lobby
	GomokuType string
	WhiteQueue *list.List
	BlackQueue *list.List
	PlayerMap  map[*core.Player]*GomokuLobbySlot
	LobbyType string
	mu sync.RWMutex
}

type GomokuLobbySlot struct {
	Element *list.Element
	Queue *list.List
}

func NewGomokuLobby(maxPlayers int, gomokuType string, db *db.Database) core.LobbyController {
	gomokuLobby := &GomokuLobby{
		Lobby: &core.Lobby{
			NumPlayers: 0,
			MaxPlayers: maxPlayers,
			RoomManager: core.NewRoomManager(),
			DB: db,
		},
		GomokuType: gomokuType,
		WhiteQueue: list.New(),
		BlackQueue: list.New(),
		LobbyType: gomokuType,
		PlayerMap: make(map[*core.Player]*GomokuLobbySlot),
	}
	
	go gomokuLobby.MatchPlayers()

	return gomokuLobby
}

func (lobby *GomokuLobby) AddPlayer(player *core.Player) {
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

func (lobby *GomokuLobby) RemovePlayer(player *core.Player) {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	slot, ok := lobby.PlayerMap[player]
	if !ok { return }

	slot.Queue.Remove(slot.Element)
	delete(lobby.PlayerMap, player)
	lobby.NumPlayers--
}


func (lobby *GomokuLobby) MatchPlayers() {
	for {
		w, b, ok := lobby.tryMatch()
		if ok {
			log.Println("Matched:", w.PlayerID, b.PlayerID)

			room := NewGomokuRoom(w, b, lobby.LobbyType, lobby.DB)
			if room != nil {
				lobby.RoomManager.RegisterPlayerToRoom(w.PlayerID, room)
				lobby.RoomManager.RegisterPlayerToRoom(b.PlayerID, room)
				w.StartPlayer()
				b.StartPlayer()
				room.Start()
			}
		}

		// Nothing to match — avoid busy spinning
		time.Sleep(100 * time.Millisecond)
	}
}


//Private logic
func (lobby *GomokuLobby) tryMatch() (*core.Player, *core.Player, bool) {
    lobby.mu.Lock()
    defer lobby.mu.Unlock()

    for lobby.WhiteQueue.Len() > 0 && lobby.BlackQueue.Len() > 0 {

        w := lobby.WhiteQueue.Front().Value.(*core.Player)
        b := lobby.BlackQueue.Front().Value.(*core.Player)

        if w.Conn == nil {
            lobby.WhiteQueue.Remove(lobby.WhiteQueue.Front())
            delete(lobby.PlayerMap, w)
            lobby.NumPlayers--
            continue
        }
        if b.Conn == nil {
            lobby.BlackQueue.Remove(lobby.BlackQueue.Front())
            delete(lobby.PlayerMap, b)
            lobby.NumPlayers--
            continue
        }
        if w.PlayerID == b.PlayerID {
            // remove both if corrupted
            lobby.WhiteQueue.Remove(lobby.WhiteQueue.Front())
            lobby.BlackQueue.Remove(lobby.BlackQueue.Front())
            delete(lobby.PlayerMap, w)
            delete(lobby.PlayerMap, b)
            lobby.NumPlayers -= 2
            continue
        }

        // VALID MATCH — remove from queues here under lock
        lobby.WhiteQueue.Remove(lobby.WhiteQueue.Front())
        lobby.BlackQueue.Remove(lobby.BlackQueue.Front())
        delete(lobby.PlayerMap, w)
        delete(lobby.PlayerMap, b)
        lobby.NumPlayers -= 2

        return w, b, true
    }

    return nil, nil, false
}
