package game

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)


type Player struct {
	PlayerID       string `json:"playerID"`
	Color string `json:"color"`
	PlayerName		string `json:"playerName"`
	Clock *PlayerClock `json:"playerClock"`
	Conn     *websocket.Conn `json:"-"`
  Incoming chan []byte `json:"-"`
  Outgoing chan []byte `json:"-"`
	Disconnected atomic.Bool `json:"-"`
	closeOnce sync.Once `json:"-"`
}

type PlayerClock struct {
	Remaining time.Duration `json:"remaining"`
	IsActive atomic.Bool `json:"-"`
	Timeout chan string `json:"-"`
}

func NewPlayers(p1 *Player, p2 *Player) []*Player {
	newPlayers := make([]*Player, 2)
	newPlayers[0] = p1
	newPlayers[1] = p2

	return newPlayers
}

func (player* Player) StartPlayer() {
	player.StartReader()
	player.StartWriter()
	if player.Color == "black" {
		player.StartClock()
	}
}

func (player *Player) StartReader() {
	go func() {
		for {
			if player.Disconnected.Load() { continue }

			_, msg, err := player.Conn.ReadMessage()
			if err != nil {
				player.Disconnected.Store(true)

			} else {
				player.Disconnected.Store(false)
				player.Incoming <- msg
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
			} else {
				player.Disconnected.Store(false)
			}
		}
	}()
}


func (player *Player) StartClock() {
	ticker := time.NewTicker(1 * time.Second)

	player.Clock.IsActive.Store(true)
	lastTick := time.Now()
	
	
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if !player.Clock.IsActive.Load() {
				return
			}
				elapsed := time.Since(lastTick)
				player.Clock.Remaining -= elapsed
				lastTick = time.Now()
			
			if player.Clock.Remaining <= 0 {
				player.StopClock()
				select {
					case player.Clock.Timeout <- player.PlayerID:
					default:
				}
				return
			}
		}
	}()
}

func (player *Player) StopClock() {
	player.Clock.IsActive.Store(false)
}

func (player *Player) ClosePlayer() {
	player.closeOnce.Do(func() {
		close(player.Incoming)
		close(player.Outgoing)
		player.Conn.Close()
	})
}

func GetPlayerByColor(gameState *GameState, color string) *Player {
	if gameState.Players[0].Color == color {
		return gameState.Players[0]
	}
	return gameState.Players[1]
}

func GetPlayerByID(gameState *GameState, playerID string) *Player {
	if gameState.Players[0].PlayerID == playerID {
		return gameState.Players[0]
	}
	return gameState.Players[1]
}

func GetOpponentPlayerByColor(gameState *GameState, color string) *Player {
	if gameState.Players[0].Color != color {
		return gameState.Players[0]
	}
	return gameState.Players[1]
}
