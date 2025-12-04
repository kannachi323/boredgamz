package core

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerController interface {
	StartPlayer()
	ClosePlayer()
	StartReader()
	StartWriter()
	StartClock()
	TimeoutClock()
}

type Player struct {
	PlayerID    string         `json:"playerID"`
	Color       string         `json:"color"`
	PlayerName  string         `json:"playerName"`
	Clock       *PlayerClock   `json:"playerClock"`
	Conn        *websocket.Conn `json:"-"`
	Incoming    chan []byte     `json:"-"`
	Outgoing    chan []byte     `json:"-"`
	Disconnected atomic.Bool     `json:"-"`
	CloseOnce   sync.Once        `json:"-"`
}

type PlayerClock struct {
	Remaining time.Duration `json:"remaining"`
	IsActive  atomic.Bool   `json:"-"`
	Timeout      chan []byte `json:"-"`
}

func NewPlayerClock(remaining time.Duration) *PlayerClock {
	return &PlayerClock{
		Remaining: remaining,
		Timeout:  make(chan []byte, 10),
	}
}

func NewPlayer(playerID, playerName, color string, clock *PlayerClock, conn *websocket.Conn) *Player {
	return &Player{
		PlayerID:    playerID,
		PlayerName:  playerName,
		Color:       color,
		Clock:       clock,
		Conn:        conn,
		Incoming:    make(chan []byte, 10),
		Outgoing:    make(chan []byte, 10),
		Disconnected: atomic.Bool{},
		CloseOnce:   sync.Once{},
	}
}

func (p *Player) StartPlayer() {
	p.StartReader()
	p.StartWriter()
	p.RunClock()
}

func (p *Player) ClosePlayer() {
	p.CloseOnce.Do(func() {
		p.Disconnected.Store(true)
		close(p.Clock.Timeout)
		close(p.Incoming)
		close(p.Outgoing)
		p.Conn.Close()
	})
}

func (p *Player) StartReader() {
	go func() {
		for {
			if p.Disconnected.Load() {
				return
			}

			_, msg, err := p.Conn.ReadMessage()
			if err != nil {
				p.Disconnected.Store(true)
				return
			}

			select {
			case p.Incoming <- msg:
			default:
				log.Println("Player incoming channel full - dropping message")
			}
		}
	}()
}

func (p *Player) StartWriter() {
	go func() {
		for msg := range p.Outgoing {
			if p.Disconnected.Load() {
				return
			}

			err := p.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				p.Disconnected.Store(true)
				return
			}
		}
	}()
}


func (p *Player) StartClock() {
	p.Clock.IsActive.Store(true)
}

func (p *Player) StopClock() {
	p.Clock.IsActive.Store(false)
}

func (p *Player) RunClock() {
	ticker := time.NewTicker(time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if p.Clock.IsActive.Load() {
					p.Clock.Remaining -= time.Second

					if p.Clock.Remaining <= 0 {
						p.Clock.Remaining = 0
						p.Clock.IsActive.Store(false)
						select {
						case p.Clock.Timeout <- []byte(p.PlayerID):
						}
					}
				}

			case <-p.Clock.Timeout:
				return
			}
		}
	}()
}
