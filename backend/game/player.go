package game

import (
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Conn     *websocket.Conn
    Incoming chan []byte
    Outgoing chan []byte
}

func NewPlayers(p1 Player, p2 Player) []*Player {
	newPlayers := make([]*Player, 2)
	newPlayers[0] = &p1
	newPlayers[1] = &p2

	return newPlayers
}

func (player *Player) StartReader() {
	go func() {
		for {
			_, msg, err := player.Conn.ReadMessage()
			if err != nil {
				log.Printf("closing room...")
				close(player.Incoming)
				return
			}
			player.Incoming <- msg
		}
	}()
}

func (player *Player) StartWriter() {
	go func() {
		for msg := range player.Outgoing {
			err := player.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("cannot write message...")
				return
			}
		}
	}()
}