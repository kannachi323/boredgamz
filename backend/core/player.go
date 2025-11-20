package core

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerController interface {
	//Player lifecycle methods
	StartPlayer()
	ClosePlayer()
	StartReader()
	StartWriter()
	StartClock()
	StopClock()

}

type Player struct {
	PlayerID       string `json:"playerID"`
	Color string `json:"color"`
	PlayerName		string `json:"playerName"`
	Clock *PlayerClock `json:"playerClock"`
	Conn     *websocket.Conn `json:"-"`
	Incoming chan []byte `json:"-"`
	Outgoing chan []byte `json:"-"`
	Disconnected atomic.Bool `json:"-"`
	CloseOnce sync.Once `json:"-"`
}

type PlayerClock struct {
	Remaining time.Duration `json:"remaining"`
	IsActive atomic.Bool `json:"-"`
	Timeout chan []byte `json:"-"`
}

func NewPlayer(playerID, playerName, color string, clock *PlayerClock, conn *websocket.Conn) *Player{
	return &Player{
		PlayerID:      playerID,
		PlayerName:    playerName,
		Color:         color,
		Clock:      clock,
		Conn:     conn,
		Incoming: make(chan []byte, 10),
		Outgoing: make(chan []byte, 10),
		Disconnected: atomic.Bool{},
		CloseOnce: sync.Once{},
	}
}


func (player *Player) StartPlayer() {
	player.StartReader()
	player.StartWriter()
	player.RunClock()
}

func (player *Player) ClosePlayer() {
	player.CloseOnce.Do(func() {
		close(player.Incoming)
		close(player.Outgoing)
		player.Conn.Close()
	})
}


func (player *Player) StartReader() {
	go func() {
		for {
			if player.Disconnected.Load() { continue }
			
			_, msg, err := player.Conn.ReadMessage()
            if err != nil {
                player.Disconnected.Store(true)
                return 
            }
			select {
            case player.Incoming <- msg:
            default:
                log.Println("Player incoming channel full - dropping message")
            }
		}
	}()
}

func (player *Player) StartWriter() {
	go func() {
		for msg := range player.Outgoing {
			if player.Disconnected.Load() { continue }
			
			err := player.Conn.WriteMessage(websocket.TextMessage, msg)
            
            if err != nil {
                player.Disconnected.Store(true)
                return
            }
        }
	}()
}


func (player *Player) StartClock() {
    if player.Clock.IsActive.Load() {
        return
    }

    player.Clock.IsActive.Store(true)
}

func (player *Player) StopClock() {
    player.Clock.IsActive.Store(false)
}

func (player *Player) RunClock() {
    ticker := time.NewTicker(time.Second)
    lastTick := time.Now()

    go func() {
        for range ticker.C {
            if player.Clock.IsActive.Load() {
                now := time.Now()
                elapsed := now.Sub(lastTick)
                player.Clock.Remaining -= elapsed
                lastTick = now

                if player.Clock.Remaining <= 0 {
                    player.Clock.IsActive.Store(false)
                }
            } else {
                lastTick = time.Now()
            }
        }
    }()
}




