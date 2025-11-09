package game

import (
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)


type Player struct {
	PlayerID       string `json:"playerID"`
	Color string `json:"color"`
	PlayerName		string `json:"playerName"`
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