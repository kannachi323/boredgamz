package connectfour

import (
	"boredgamz/core"
	"boredgamz/db"
	"container/list"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const LobbyID = "connectfour"

type ConnectFourLobby struct {
	*core.Lobby
	RedQueue    *list.List
	YellowQueue *list.List
	PlayerSlot  map[string]*ConnectFourLobbySlot

	mu     sync.Mutex
	wakeup chan struct{}
}

type ConnectFourLobbySlot struct {
	Element *list.Element
	Queue   *list.List
	Player  *core.Player
}

func NewConnectFourLobby(maxPlayers int, db *db.Database) core.LobbyController {
	lobby := &ConnectFourLobby{
		Lobby: &core.Lobby{
			LobbyName:   LobbyID,
			NumPlayers:  0,
			MaxPlayers:  maxPlayers,
			RoomManager: core.NewRoomManager(),
			DB:          db,
		},
		RedQueue:    list.New(),
		YellowQueue: list.New(),
		PlayerSlot:  make(map[string]*ConnectFourLobbySlot),
		wakeup:      make(chan struct{}, 1),
	}

	go lobby.MatchPlayers()
	return lobby
}

func (lobby *ConnectFourLobby) AddPlayer(player *core.Player) {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	if !isPlayerConnected(player) {
		return
	}

	lobby.cleanupDisconnectedLocked()

	if lobby.NumPlayers >= lobby.MaxPlayers {
		return
	}

	if slot, exists := lobby.PlayerSlot[player.PlayerID]; exists {
		if slot.Player != nil {
			slot.Player.ReconnectPlayer(player.Conn)
			lobby.attachCloseHandler(slot.Player)
		}
		return
	}

	var elem *list.Element
	var queue *list.List
	if player.Color == "" || player.Color == "random" {
		if lobby.RedQueue.Len() <= lobby.YellowQueue.Len() {
			player.Color = "red"
		} else {
			player.Color = "yellow"
		}
	}

	switch player.Color {
	case "red":
		elem = lobby.RedQueue.PushBack(player)
		queue = lobby.RedQueue
	case "yellow":
		elem = lobby.YellowQueue.PushBack(player)
		queue = lobby.YellowQueue
	default:
		return
	}

	lobby.PlayerSlot[player.PlayerID] = &ConnectFourLobbySlot{
		Element: elem,
		Queue:   queue,
		Player:  player,
	}

	player.StartPlayer()
	lobby.attachCloseHandler(player)
	lobby.NumPlayers++

	select {
	case lobby.wakeup <- struct{}{}:
	default:
	}
}

func (lobby *ConnectFourLobby) RemovePlayer(player *core.Player) {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	lobby.removePlayerByIDLocked(player.PlayerID)
}

func (lobby *ConnectFourLobby) removePlayerByIDLocked(playerID string) {
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
}

func (lobby *ConnectFourLobby) MatchPlayers() {
	for {
		<-lobby.wakeup
		for {
			red, yellow, ok := lobby.tryMatch()
			if !ok {
				break
			}

			room := NewConnectFourRoom(red, yellow, LobbyID, lobby.DB, lobby.RoomManager)
			lobby.RoomManager.RegisterPlayerToRoom(red.PlayerID, room)
			lobby.RoomManager.RegisterPlayerToRoom(yellow.PlayerID, room)
			room.Start()
		}
	}
}

func (lobby *ConnectFourLobby) tryMatch() (*core.Player, *core.Player, bool) {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()

	lobby.cleanupDisconnectedLocked()

	for lobby.RedQueue.Len() > 0 && lobby.YellowQueue.Len() > 0 {
		rElem := lobby.RedQueue.Front()
		yElem := lobby.YellowQueue.Front()
		if rElem == nil || yElem == nil {
			break
		}

		red := rElem.Value.(*core.Player)
		yellow := yElem.Value.(*core.Player)

		if !isPlayerConnected(red) {
			lobby.removePlayerByIDLocked(red.PlayerID)
			continue
		}
		if !isPlayerConnected(yellow) {
			lobby.removePlayerByIDLocked(yellow.PlayerID)
			continue
		}
		if red.PlayerID == yellow.PlayerID {
			lobby.removePlayerByIDLocked(red.PlayerID)
			lobby.removePlayerByIDLocked(yellow.PlayerID)
			continue
		}

		lobby.removePlayerByIDLocked(red.PlayerID)
		lobby.removePlayerByIDLocked(yellow.PlayerID)
		return red, yellow, true
	}

	return nil, nil, false
}

func (lobby *ConnectFourLobby) cleanupDisconnectedLocked() {
	for playerID, slot := range lobby.PlayerSlot {
		if slot == nil || slot.Player == nil || !isPlayerConnected(slot.Player) {
			lobby.removePlayerByIDLocked(playerID)
		}
	}
}

func isPlayerConnected(player *core.Player) bool {
	if player == nil || player.Conn == nil {
		return false
	}
	return !player.Disconnected.Load()
}

func (lobby *ConnectFourLobby) attachCloseHandler(player *core.Player) {
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
