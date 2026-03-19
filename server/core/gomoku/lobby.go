package gomoku

import (
	"boredgamz/core"
	"boredgamz/db"
	"container/list"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type GomokuLobby struct {
	*core.Lobby
	WhiteQueue *list.List
	BlackQueue *list.List
	PlayerSlot map[string]*GomokuLobbySlot

	Mu     sync.Mutex
	wakeup chan struct{}
}

type GomokuLobbySlot struct {
	Element *list.Element
	Queue   *list.List
	Player  *core.Player
}

func NewGomokuLobby(maxPlayers int, name string, db *db.Database) core.LobbyController {
	gomokuLobby := &GomokuLobby{
		Lobby: &core.Lobby{
			LobbyName: name,
			NumPlayers: 0,
			MaxPlayers: maxPlayers,
			RoomManager: core.NewRoomManager(),
			DB: db,
		},
		WhiteQueue: list.New(),
		BlackQueue: list.New(),
		PlayerSlot: make(map[string]*GomokuLobbySlot),
		wakeup:     make(chan struct{}, 1),
	}

	// start matcher goroutine
	go gomokuLobby.MatchPlayers()

	return gomokuLobby
}

func (lobby *GomokuLobby) AddPlayer(player *core.Player) {
	lobby.Mu.Lock()
	defer lobby.Mu.Unlock()

	if !isPlayerConnected(player) {
		log.Println("Player disconnected, not adding to queue:", player.PlayerID)
		return
	}

	lobby.cleanupDisconnectedLocked()

	if lobby.NumPlayers >= lobby.MaxPlayers {
		return
	}
	if slot, exists := lobby.PlayerSlot[player.PlayerID]; exists {
		if slot.Player != nil {
			log.Println("Player already queued, refreshing session:", player.PlayerID)
			slot.Player.ReconnectPlayer(player.Conn)
			lobby.attachCloseHandler(slot.Player)
		}
		return
	}

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
	lobby.PlayerSlot[player.PlayerID] = &GomokuLobbySlot{
		Element: elem,
		Queue:   queue,
		Player:  player,
	}

	player.StartPlayer()
	lobby.attachCloseHandler(player)

	lobby.NumPlayers++

	log.Println(lobby.WhiteQueue)
	log.Println(lobby.BlackQueue)

	select {
	case lobby.wakeup <- struct{}{}:
	default:
	}
}

func (lobby *GomokuLobby) RemovePlayer(player *core.Player) {
	lobby.Mu.Lock()
	defer lobby.Mu.Unlock()

	lobby.removePlayerByIDLocked(player.PlayerID)
}

func (lobby *GomokuLobby) removePlayerByIDLocked(playerID string) {
	slot, ok := lobby.PlayerSlot[playerID]
	if !ok {
		return
	}

	if slot.Element != nil && slot.Queue != nil {
		slot.Queue.Remove(slot.Element)
	}
	delete(lobby.PlayerSlot, playerID)
	
	if lobby.NumPlayers > 0 {
		lobby.NumPlayers--
	}
	log.Println("removed " + playerID)
	select {
	case lobby.wakeup <- struct{}{}:
	default:
	}
}

func (lobby *GomokuLobby) MatchPlayers() {
	for {
		<-lobby.wakeup
		for {
			w, b, ok := lobby.tryMatch()
			if !ok {
				break
			}

			log.Println("Matched:", w.PlayerID, b.PlayerID)

			go func(wp, bp *core.Player) {
				room := NewGomokuRoom(wp, bp, lobby.LobbyName, lobby.DB, lobby.RoomManager)
				if room == nil {
					return
				}
				lobby.RoomManager.RegisterPlayerToRoom(wp.PlayerID, room)
				lobby.RoomManager.RegisterPlayerToRoom(bp.PlayerID, room)
				room.Start()
			}(w, b)
		}
	}
}


func (lobby *GomokuLobby) tryMatch() (*core.Player, *core.Player, bool) {
	lobby.Mu.Lock()
	defer lobby.Mu.Unlock()

	lobby.cleanupDisconnectedLocked()

    for lobby.WhiteQueue.Len() > 0 && lobby.BlackQueue.Len() > 0 {
        wElem := lobby.WhiteQueue.Front()
        bElem := lobby.BlackQueue.Front()
		if wElem == nil || bElem == nil {
			continue
		}

        w := wElem.Value.(*core.Player)
        b := bElem.Value.(*core.Player)

		if !isPlayerConnected(w) {
			lobby.removePlayerByIDLocked(w.PlayerID)
			continue
		}
		if !isPlayerConnected(b) {
			lobby.removePlayerByIDLocked(b.PlayerID)
			continue
		}
		if w.PlayerID == b.PlayerID {
			lobby.removePlayerByIDLocked(w.PlayerID)
			lobby.removePlayerByIDLocked(b.PlayerID)
			continue
		}

		lobby.removePlayerByIDLocked(w.PlayerID)
		lobby.removePlayerByIDLocked(b.PlayerID)

        return w, b, true
    }

    return nil, nil, false
}

func isPlayerConnected(player *core.Player) bool {
	if player == nil || player.Conn == nil {
		return false
	}

	return !player.Disconnected.Load()
}

func (lobby *GomokuLobby) cleanupDisconnectedLocked() {
	for playerID, slot := range lobby.PlayerSlot {
		if slot == nil || slot.Player == nil {
			lobby.removePlayerByIDLocked(playerID)
			continue
		}

		if !isPlayerConnected(slot.Player) {
			lobby.removePlayerByIDLocked(playerID)
		}
	}
}

func (lobby *GomokuLobby) attachCloseHandler(player *core.Player) {
	if player == nil || player.Conn == nil {
		return
	}

	player.Conn.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, text)
		_ = player.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		lobby.RemovePlayer(player)
		return nil
	})
}
