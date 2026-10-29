package arena

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 4096
)

func ServePlayer(hub *Hub, options JoinOptions, conn *websocket.Conn) error {
	player, room, err := hub.JoinMatch(options)
	if err != nil {
		return err
	}
	generation := player.generation.Load()
	send := player.sender()
	defer func() {
		if room == nil {
			room = hub.assignedRoom(player)
		}
		if room == nil {
			player.connected.Store(false)
			hub.LeaveWaiting(player)
		} else {
			hub.Disconnect(player, room, generation)
		}
		conn.Close()
	}()

	roomReady := make(chan *Room, 1)
	if room == nil {
		send <- ServerMessage{Type: "waiting", Message: "Looking for an opponent"}
		go func(ready chan<- *Room) { ready <- hub.MatchWaiting(player) }(roomReady)
	} else {
		room.Broadcast()
	}
	readDone := make(chan struct{})
	incoming := make(chan ClientMessage, 8)
	go readMessages(conn, incoming, readDone)

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case message := <-send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(message); err != nil {
				return err
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		case <-readDone:
			return nil
		case matched := <-roomReady:
			if matched == nil {
				return nil
			}
			// The opponent's join already broadcast the first state to us.
			room = matched
			roomReady = nil
		case message := <-incoming:
			if room == nil {
				// Matched, but the poller hasn't told this loop yet.
				if assigned := hub.assignedRoom(player); assigned != nil {
					room, roomReady = assigned, nil
				}
			}
			if room == nil {
				sendError(player, "Wait for an opponent before moving")
				continue
			}
			if message.Type == "leave" {
				hub.LeaveRoom(player, room)
				return nil
			}
			if err := room.HandleClient(player.ID, message); err != nil {
				sendError(player, err.Error())
			}
		}
	}
}

func readMessages(conn *websocket.Conn, incoming chan<- ClientMessage, done chan<- struct{}) {
	defer close(done)
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var message ClientMessage
		if err := json.Unmarshal(payload, &message); err != nil {
			continue
		}
		select {
		case incoming <- message:
		default:
		}
	}
}

func sendError(player *Player, message string) {
	send := player.sender()
	select {
	case send <- ServerMessage{Type: "error", Message: message}:
	default:
		log.Printf("dropping arena error for player %s", player.ID)
	}
}
