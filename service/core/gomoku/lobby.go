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
	PlayerSlot  map[*core.Player]*GomokuLobbySlot

	Mu     sync.Mutex
	wakeup chan struct{}           
}

type GomokuLobbySlot struct {
	Element *list.Element
	Queue   *list.List
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
		PlayerSlot: make(map[*core.Player]*GomokuLobbySlot),
		wakeup:    make(chan struct{}, 1),
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

	if lobby.NumPlayers >= lobby.MaxPlayers {
		return
	}
	if _, exists := lobby.PlayerSlot[player]; exists {
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
	lobby.PlayerSlot[player] = &GomokuLobbySlot{
		Element: elem,
		Queue:   queue,
	}

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

	lobby.removePlayer(player)
}

func (lobby* GomokuLobby) removePlayer(player *core.Player) {
	slot, ok := lobby.PlayerSlot[player]
	if !ok {
		return
	}

	if slot.Element != nil && slot.Queue != nil {
		slot.Queue.Remove(slot.Element)
	}
	delete(lobby.PlayerSlot, player)
	
	if lobby.NumPlayers > 0 {
		lobby.NumPlayers--
	}
	log.Println("removed " + player.PlayerID)
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
				room := NewGomokuRoom(wp, bp, lobby.LobbyName, lobby.DB)
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

    for lobby.WhiteQueue.Len() > 0 && lobby.BlackQueue.Len() > 0 {
        wElem := lobby.WhiteQueue.Front()
        bElem := lobby.BlackQueue.Front()
		if (wElem == nil || bElem == nil) { continue }

        w := wElem.Value.(*core.Player)
        b := bElem.Value.(*core.Player)

		if !isPlayerConnected(w) {
			lobby.removePlayer(w)
			continue
		}
		if !isPlayerConnected(b) {
			lobby.removePlayer(b)
			continue
		}
		if w.PlayerID == b.PlayerID {
			lobby.removePlayer(w)
			lobby.removePlayer(b)
			continue
		}

		lobby.removePlayer(w)
		lobby.removePlayer(b)

        return w, b, true
    }

    return nil, nil, false
}

func isPlayerConnected(player *core.Player) bool {
    err := player.Conn.WriteControl(
        websocket.PingMessage,
        []byte{},
        time.Now().Add(time.Second),
    )
    return err == nil
}
