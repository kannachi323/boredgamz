package core

import (
	"context"
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
	PlayerID     string `json:"playerID"`
	Color        string `json:"color"`
	PlayerName   string `json:"playerName"`
	Clock        *PlayerClock `json:"playerClock"`
	Conn         *websocket.Conn `json:"-"`
	Incoming     chan []byte `json:"-"`
	Outgoing     chan []byte `json:"-"`
	Disconnected atomic.Bool `json:"-"`
	
	// Context manages the lifecycle of the Writer/Reader loops
	Ctx          context.Context    `json:"-"`
	Cancel       context.CancelFunc `json:"-"`
	
	// Use a pointer for Once so we can safely replace it on reconnect
	CloseOnce    *sync.Once         `json:"-"`
}

type PlayerClock struct {
	Remaining time.Duration `json:"remaining"`
}

func NewPlayerClock(timeControl string) *PlayerClock {
	var remaining time.Duration
	// Using time.Second is cleaner and safer than raw nanoseconds
	switch timeControl {
	case "Rapid":
		remaining = 300 * time.Second // 5 mins
	case "Blitz":
		remaining = 180 * time.Second // 3 mins
	case "Bullet":
		remaining = 60 * time.Second  // 1 min
	case "Hyperbullet":
		remaining = 30 * time.Second  // 30 secs
	default:
		remaining = 300 * time.Second
	}
	return &PlayerClock{
		Remaining: remaining,
	}
}

func NewPlayer(playerID, playerName, color string, clock *PlayerClock, conn *websocket.Conn) *Player {
	ctx, cancel := context.WithCancel(context.Background())
	return &Player{
		PlayerID:     playerID,
		PlayerName:   playerName,
		Color:        color,
		Clock:        clock,
		Conn:         conn,
		Incoming:     make(chan []byte, 10),
		Outgoing:     make(chan []byte, 10),
		Disconnected: atomic.Bool{},
		CloseOnce:    &sync.Once{},
		Ctx:          ctx,
		Cancel:       cancel,
	}
}

func (p *Player) ReconnectPlayer(conn *websocket.Conn) {
	// 1. KILL THE OLD LOOPS FIRST
	// This stops the "Old Writer" from stealing messages destined for the new connection
	p.Cancel()

	// 2. Create new context for the new session
	p.Ctx, p.Cancel = context.WithCancel(context.Background())

	// 3. Reset state
	p.Disconnected.Store(false)
	p.CloseOnce = &sync.Once{} // Fresh sync.Once
	p.Conn = conn

	// 4. Start new loops
	p.StartReader()
	p.StartWriter()
}

func (p *Player) StartPlayer() {
	p.StartReader()
	p.StartWriter()
}

func (p *Player) ClosePlayer() {
	if p.CloseOnce == nil {
		return
	}
	p.CloseOnce.Do(func() {
		p.Disconnected.Store(true)
		p.Conn.Close()
		p.Cancel() // Ensure all loops stop immediately
	})
}

func (p *Player) StartReader() {
	go func() {
		// Ensure we clean up if this loop dies
		defer p.ClosePlayer()

		for {
			select {
			case <-p.Ctx.Done():
				return // Stop if player was reconnected/closed
			default:
				_, msg, err := p.Conn.ReadMessage()
				if err != nil {
					log.Println("Read error (player disconnected):", err)
					return // defer will call ClosePlayer()
				}

				select {
				case p.Incoming <- msg:
				default:
					log.Println("Player incoming channel full - dropping message")
				}
			}
		}
	}()
}

func (p *Player) StartWriter() {
	// Capture the connection valid for THIS loop instance
	conn := p.Conn 
	
	go func() {
		for {
			select {
			// 1. Stop if the context is cancelled (Reconnected or Closed)
			case <-p.Ctx.Done():
				return

			// 2. Otherwise, process outgoing messages
			case msg := <-p.Outgoing:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Write error:", err)
					p.ClosePlayer() // Signal that the connection is dead
					return 
				}
			}
		}
	}()
}