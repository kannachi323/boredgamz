package poker

import (
	"errors"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Room struct {
	id             string
	table          *Table
	players        [2]*Player
	spectators     map[string]*Player
	inviteCode     string
	chat           []ChatMessage
	lastChat       map[string]time.Time
	ready          map[int]bool
	actionTimer    *time.Timer
	timerVersion   uint64
	actionDeadline *time.Time
	mu             sync.Mutex
}

func newRoom(id, mode string, first, second *Player, invite string) (*Room, error) {
	first.Seat, second.Seat = 0, 1
	table, err := NewTable(mode, &Seat{ID: first.ID, Name: first.Name}, &Seat{ID: second.ID, Name: second.Name})
	if err != nil {
		return nil, err
	}
	room := &Room{id: id, table: table, players: [2]*Player{first, second}, spectators: map[string]*Player{}, inviteCode: invite, lastChat: map[string]time.Time{}, ready: map[int]bool{}}
	room.resetActionTimerLocked()
	return room, nil
}

func (r *Room) Handle(playerID string, message ClientMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	player := r.participantLocked(playerID)
	if player == nil {
		return errors.New("player is not in this poker room")
	}
	switch message.Type {
	case "action":
		if player.Seat < 0 {
			return errors.New("spectators cannot wager")
		}
		if err := r.table.Act(player.Seat, message.Action, message.Amount); err != nil {
			return err
		}
		r.resetActionTimerLocked()
	case "discard":
		if player.Seat < 0 {
			return errors.New("spectators cannot draw")
		}
		if err := r.table.Discard(player.Seat, message.Positions); err != nil {
			return err
		}
		r.resetActionTimerLocked()
	case "chat":
		if err := r.addChatLocked(player, message.Content); err != nil {
			return err
		}
	case "next-hand":
		if player.Seat < 0 {
			return errors.New("spectators cannot ready")
		}
		if r.table.Phase != "showdown" && r.table.Phase != "hand-finished" {
			return errors.New("the hand is still active")
		}
		r.ready[player.Seat] = true
		if r.ready[0] && r.ready[1] {
			if r.table.Seats[0].Stack == 0 || r.table.Seats[1].Stack == 0 {
				r.table.Seats[0].Stack, r.table.Seats[1].Stack = StartingStack, StartingStack
			}
			if err := r.table.StartHand(); err != nil {
				return err
			}
			r.ready = map[int]bool{}
			r.resetActionTimerLocked()
		}
	default:
		return errors.New("unsupported poker message")
	}
	r.broadcastLocked()
	return nil
}

func (r *Room) addChatLocked(player *Player, content string) error {
	content = strings.TrimSpace(content)
	chars := []rune(content)
	if len(chars) > 240 {
		chars = chars[:240]
	}
	content = strings.Map(func(value rune) rune {
		if unicode.IsControl(value) && value != '\n' && value != '\t' {
			return -1
		}
		return value
	}, string(chars))
	if content == "" {
		return errors.New("message is empty")
	}
	if time.Since(r.lastChat[player.ID]) < 750*time.Millisecond {
		return errors.New("messages are being sent too quickly")
	}
	r.lastChat[player.ID] = time.Now()
	r.chat = append(r.chat, ChatMessage{ID: uuid.NewString(), SenderID: player.ID, SenderName: player.Name, Content: content, CreatedAt: time.Now().UTC()})
	if len(r.chat) > 60 {
		r.chat = append([]ChatMessage{}, r.chat[len(r.chat)-60:]...)
	}
	return nil
}

func (r *Room) snapshotLocked(viewer *Player) *Snapshot {
	s := &Snapshot{RoomID: r.id, Mode: r.table.Mode, Phase: r.table.Phase, Community: append([]Card{}, r.table.Community...), Dealer: r.table.Dealer, Turn: r.table.Turn, Pot: r.table.Pot, AwardedPot: r.table.AwardedPot, CurrentBet: r.table.CurrentBet, MinimumRaise: r.table.MinimumRaise, Hand: r.table.Hand, Winners: append([]int{}, r.table.Winners...), WinningHand: r.table.WinningHand, Message: r.table.Message, Chat: append([]ChatMessage{}, r.chat...), Spectators: len(r.spectators), InviteCode: r.inviteCode, ActionDeadline: clonePokerTime(r.actionDeadline)}
	for index, seat := range r.table.Seats {
		view := &SeatView{ID: seat.ID, Name: seat.Name, Stack: seat.Stack, RoundBet: seat.RoundBet, TotalBet: seat.TotalBet, Folded: seat.Folded, AllIn: seat.AllIn, Connected: r.players[index].connected.Load(), CardCount: len(seat.Cards)}
		if r.table.Phase == "showdown" || (viewer != nil && viewer.Seat == index) {
			view.Cards = append([]Card{}, seat.Cards...)
		}
		s.Seats[index] = view
	}
	if viewer != nil && viewer.Seat >= 0 {
		s.ToCall = r.table.ToCall(viewer.Seat)
		s.LegalActions = r.table.LegalActions(viewer.Seat)
	}
	for index, value := range r.ready {
		if value {
			s.Ready = append(s.Ready, index)
		}
	}
	if s.Chat == nil {
		s.Chat = []ChatMessage{}
	}
	if s.Winners == nil {
		s.Winners = []int{}
	}
	if s.Ready == nil {
		s.Ready = []int{}
	}
	if s.Community == nil {
		s.Community = []Card{}
	}
	if s.LegalActions == nil {
		s.LegalActions = []string{}
	}
	return s
}

func (r *Room) broadcastLocked() {
	for _, player := range r.players {
		r.sendLocked(player)
	}
	for _, player := range r.spectators {
		r.sendLocked(player)
	}
}
func (r *Room) sendLocked(player *Player) {
	if player == nil || !player.connected.Load() {
		return
	}
	select {
	case player.sender() <- ServerMessage{Type: "state", State: r.snapshotLocked(player)}:
	default:
	}
}
func (r *Room) SendState(player *Player) { r.mu.Lock(); defer r.mu.Unlock(); r.sendLocked(player) }
func (r *Room) participantLocked(id string) *Player {
	for _, p := range r.players {
		if p != nil && p.ID == id {
			return p
		}
	}
	return r.spectators[id]
}

func (r *Room) resetActionTimerLocked() {
	if r.actionTimer != nil {
		r.actionTimer.Stop()
		r.actionTimer = nil
	}
	r.actionDeadline = nil
	r.timerVersion++
	if !r.table.bettingPhase() && r.table.Phase != "draw" {
		return
	}
	deadline := time.Now().UTC().Add(30 * time.Second)
	r.actionDeadline = &deadline
	version, hand, phase, turn := r.timerVersion, r.table.Hand, r.table.Phase, r.table.Turn
	r.actionTimer = time.AfterFunc(30*time.Second, func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.timerVersion != version || r.table.Hand != hand || r.table.Phase != phase || r.table.Turn != turn {
			return
		}
		if r.table.Phase == "draw" {
			_ = r.table.Discard(turn, nil)
		} else if r.table.ToCall(turn) == 0 {
			_ = r.table.Act(turn, "check", 0)
		} else {
			_ = r.table.Act(turn, "fold", 0)
		}
		r.resetActionTimerLocked()
		r.broadcastLocked()
	})
}

func clonePokerTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
