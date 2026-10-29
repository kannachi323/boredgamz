package mahjong

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
)

type Player struct {
	ID         string
	Name       string
	Seat       int
	mu         sync.RWMutex
	send       chan ServerMessage
	connected  atomic.Bool
	generation atomic.Uint64
}

func newPlayer(id, name string) *Player {
	name = strings.TrimSpace(strings.Map(func(value rune) rune {
		if unicode.IsControl(value) {
			return -1
		}
		return value
	}, name))
	if name == "" {
		name = "Guest"
	}
	runes := []rune(name)
	if len(runes) > 24 {
		name = string(runes[:24])
	}
	player := &Player{ID: id, Name: name, Seat: -1}
	player.replaceSender()
	return player
}

func (p *Player) replaceSender() {
	p.mu.Lock()
	p.send = make(chan ServerMessage, 32)
	p.mu.Unlock()
	p.connected.Store(true)
	p.generation.Add(1)
}

func (p *Player) sender() chan ServerMessage {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.send
}

type JoinOptions struct{ PlayerID, Name, Mode, InviteCode, Role, RoomID string }

type ChatMessage struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"senderId"`
	SenderName string    `json:"senderName"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
}

type SeatView struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Pairs     int    `json:"pairs"`
	Remaining int    `json:"remaining"`
	Connected bool   `json:"connected"`
}

type Snapshot struct {
	RoomID         string        `json:"roomId"`
	Mode           string        `json:"mode"`
	Status         string        `json:"status"`
	Tiles          []Tile        `json:"tiles"`
	Removed        []bool        `json:"removed"`
	RemovedBy      []int         `json:"removedBy"`
	Free           []int         `json:"free"`
	Seats          [2]*SeatView  `json:"seats"`
	ViewerSeat     int           `json:"viewerSeat"`
	Winner         int           `json:"winner"`
	Message        string        `json:"message"`
	HintsRemaining int           `json:"hintsRemaining"`
	Hint           []int         `json:"hint"`
	Chat           []ChatMessage `json:"chat"`
	Spectators     int           `json:"spectators"`
	Ready          []int         `json:"ready"`
	InviteCode     string        `json:"inviteCode,omitempty"`
	StartedAt      time.Time     `json:"startedAt"`
	Deadline       *time.Time    `json:"deadline,omitempty"`
	Reshuffles     int           `json:"reshuffles"`
}

type ClientMessage struct {
	Type    string `json:"type"`
	First   int    `json:"first,omitempty"`
	Second  int    `json:"second,omitempty"`
	Content string `json:"content,omitempty"`
}

type ServerMessage struct {
	Type    string    `json:"type"`
	State   *Snapshot `json:"state,omitempty"`
	Message string    `json:"message,omitempty"`
}
