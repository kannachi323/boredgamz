package game

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)


type Player struct {
	PlayerID       string `json:"playerID"`
	Username string `json:"username"`
	Conn     *websocket.Conn `json:"-"`
    Incoming chan []byte `json:"-"`
    Outgoing chan []byte `json:"-"`
	Disconnected atomic.Bool `json:"-"`
	closeOnce sync.Once `json:"-"`
}

func NewPlayers(p1 *Player, p2 *Player) []*Player {
	newPlayers := make([]*Player, 2)
	newPlayers[0] = p1
	newPlayers[1] = p2

	return newPlayers
}

func (player *Player) StartReader() {
	go func() {
		for {
			_, msg, err := player.Conn.ReadMessage()
			if err != nil {
				log.Printf("Player %s disconnected: %v", player.PlayerID, err)
				player.Disconnected.Store(true)
				player.Close()
				break
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
				break
			}
		}
	}()
}

func (p *Player) Close() {
	p.closeOnce.Do(func() {
		close(p.Outgoing)
		close(p.Incoming)
		p.Conn.Close()
	})
}