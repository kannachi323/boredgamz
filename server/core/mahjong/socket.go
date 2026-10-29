package mahjong

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

const mahjongMaxMessageSize = 4096

func Serve(hub *Hub, options JoinOptions, conn *websocket.Conn) error {
	player, room, err := hub.Join(options)
	if err != nil {
		return err
	}
	generation := player.generation.Load()
	send := player.sender()
	defer func() {
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
		send <- ServerMessage{Type: "waiting", Message: "Looking for a mahjong partner"}
		go func() { roomReady <- hub.MatchWaiting(player) }()
	} else {
		room.SendState(player)
	}
	incoming := make(chan ClientMessage, 8)
	done := make(chan struct{})
	go readMessages(conn, incoming, done)
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
			room = matched
			roomReady = nil
			room.SendState(player)
		case message := <-incoming:
			if room == nil {
				sendError(player, "Wait for another player before acting")
				continue
			}
			if message.Type == "leave" {
				hub.Leave(player, room)
				return nil
			}
			if err := room.Handle(player.ID, message); err != nil {
				sendError(player, err.Error())
			}
		}
	}
}

func readMessages(conn *websocket.Conn, incoming chan<- ClientMessage, done chan<- struct{}) {
	defer close(done)
	conn.SetReadLimit(mahjongMaxMessageSize)
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

func sendError(player *Player, message string) {
	select {
	case player.sender() <- ServerMessage{Type: "error", Message: message}:
	default:
	}
}
