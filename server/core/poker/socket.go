package poker

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

const pokerMaxMessageSize = 4096

func Serve(hub *Hub, options JoinOptions, conn *websocket.Conn) error {
	player, room, err := hub.Join(options)
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
		_ = conn.Close()
	}()
	roomReady := make(chan *Room, 1)
	if room == nil {
		send <- ServerMessage{Type: "waiting", Message: "Looking for a poker opponent"}
		go func(ready chan<- *Room) { ready <- hub.MatchWaiting(player) }(roomReady)
	} else {
		room.Broadcast()
	}
	incoming := make(chan ClientMessage, 8)
	done := make(chan struct{})
	go readPokerMessages(conn, incoming, done)
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case message := <-send:
			_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteJSON(message); err != nil {
				return err
			}
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		case <-done:
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
				sendPokerError(player, "Wait for an opponent before acting")
				continue
			}
			if message.Type == "leave" {
				hub.Leave(player, room)
				return nil
			}
			if err := room.Handle(player.ID, message); err != nil {
				sendPokerError(player, err.Error())
			}
		}
	}
}

func readPokerMessages(conn *websocket.Conn, incoming chan<- ClientMessage, done chan<- struct{}) {
	defer close(done)
	conn.SetReadLimit(pokerMaxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error { return conn.SetReadDeadline(time.Now().Add(60 * time.Second)) })
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var message ClientMessage
		if json.Unmarshal(payload, &message) != nil {
			continue
		}
		select {
		case incoming <- message:
		default:
		}
	}
}
func sendPokerError(player *Player, message string) {
	select {
	case player.sender() <- ServerMessage{Type: "error", Message: message}:
	default:
	}
}
