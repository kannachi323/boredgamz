package poker

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
	mode       string // requested mode; the host's wins on private tables
	mu         sync.RWMutex
	send       chan ServerMessage
	connected  atomic.Bool
	left       atomic.Bool // pressed Leave; won't be back
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
	p.left.Store(false)
	p.generation.Add(1)
}
func (p *Player) sender() chan ServerMessage { p.mu.RLock(); defer p.mu.RUnlock(); return p.send }

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
	Stack     int    `json:"stack"`
	RoundBet  int    `json:"roundBet"`
	TotalBet  int    `json:"totalBet"`
	Folded    bool   `json:"folded"`
	AllIn     bool   `json:"allIn"`
	Connected bool   `json:"connected"`
	Left      bool   `json:"left,omitempty"`
	Cards     []Card `json:"cards,omitempty"`
	CardCount int    `json:"cardCount"`
}
type Snapshot struct {
	ViewerID       string        `json:"viewerId,omitempty"`
	ViewerSeat     int           `json:"viewerSeat"`
	RoomID         string        `json:"roomId"`
	Mode           string        `json:"mode"`
	Phase          string        `json:"phase"`
	Seats          [2]*SeatView  `json:"seats"`
	Community      []Card        `json:"community"`
	Dealer         int           `json:"dealer"`
	Turn           int           `json:"turn"`
	Pot            int           `json:"pot"`
	AwardedPot     int           `json:"awardedPot"`
	CurrentBet     int           `json:"currentBet"`
	MinimumRaise   int           `json:"minimumRaise"`
	ToCall         int           `json:"toCall"`
	LegalActions   []string      `json:"legalActions"`
	Hand           int           `json:"hand"`
	Winners        []int         `json:"winners"`
	WinningHand    string        `json:"winningHand,omitempty"`
	Message        string        `json:"message"`
	Chat           []ChatMessage `json:"chat"`
	Spectators     int           `json:"spectators"`
	Ready          []int         `json:"ready"`
	InviteCode     string        `json:"inviteCode,omitempty"`
	ActionDeadline *time.Time    `json:"actionDeadline,omitempty"`
}
type ClientMessage struct {
	Type      string `json:"type"`
	Action    string `json:"action,omitempty"`
	Amount    int    `json:"amount,omitempty"`
	Positions []int  `json:"positions,omitempty"`
	Content   string `json:"content,omitempty"`
}
type ServerMessage struct {
	Type    string    `json:"type"`
	State   *Snapshot `json:"state,omitempty"`
	Message string    `json:"message,omitempty"`
}
