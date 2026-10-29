package mahjong

import (
	"errors"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
)

const (
	gardenDuration = 12 * time.Minute
	duelDuration   = 8 * time.Minute
	startingHints  = 3
)

type Room struct {
	id         string
	mode       string
	boards     [2]*Board
	removed    [2][]bool
	removedBy  []int
	players    [2]*Player
	pairs      [2]int
	hints      [2]int
	hint       [2][]int
	spectators map[string]*Player
	inviteCode string
	chat       []ChatMessage
	lastChat   map[string]time.Time
	ready      map[int]bool
	status     string
	winner     int
	message    string
	startedAt  time.Time
	deadline   *time.Time
	timer      *time.Timer
	reshuffles int
	mu         sync.Mutex
}

func newRoom(id, mode string, first, second *Player, invite string) (*Room, error) {
	first.Seat, second.Seat = 0, 1
	r := &Room{id: id, mode: normalizeMode(mode), players: [2]*Player{first, second}, spectators: map[string]*Player{}, inviteCode: invite, lastChat: map[string]time.Time{}, ready: map[int]bool{}, winner: -1}
	if err := r.startRoundLocked(); err != nil {
		return nil, err
	}
	return r, nil
}

func normalizeMode(mode string) string {
	switch strings.TrimSpace(mode) {
	case "duel", "zen":
		return strings.TrimSpace(mode)
	default:
		return "garden"
	}
}

func (r *Room) startRoundLocked() error {
	if r.timer != nil {
		r.timer.Stop()
	}
	layout := "compact"
	if r.mode == "garden" {
		layout = "full"
	}
	board, err := NewBoard(layout)
	if err != nil {
		return err
	}
	r.boards[0] = board
	r.boards[1] = board
	if r.mode == "duel" {
		r.boards[1] = cloneBoard(board)
	}
	r.removed[0] = make([]bool, len(board.Tiles))
	if r.mode == "duel" {
		r.removed[1] = make([]bool, len(board.Tiles))
	} else {
		r.removed[1] = r.removed[0]
	}
	r.removedBy = make([]int, len(board.Tiles))
	for index := range r.removedBy {
		r.removedBy[index] = -1
	}
	r.pairs = [2]int{}
	r.hints = [2]int{startingHints, startingHints}
	r.hint = [2][]int{}
	r.ready = map[int]bool{}
	r.status, r.winner, r.message, r.reshuffles = "playing", -1, "Find two free matching tiles", 0
	r.startedAt = time.Now().UTC()
	r.deadline = nil
	var duration time.Duration
	if r.mode == "garden" {
		duration = gardenDuration
	} else if r.mode == "duel" {
		duration = duelDuration
	}
	if duration > 0 {
		deadline := r.startedAt.Add(duration)
		r.deadline = &deadline
		r.timer = time.AfterFunc(duration, func() {
			r.mu.Lock()
			defer r.mu.Unlock()
			if r.status != "playing" {
				return
			}
			r.finishTimedLocked()
			r.broadcastLocked()
		})
	}
	return nil
}

func cloneBoard(board *Board) *Board {
	copy := &Board{Tiles: append([]Tile{}, board.Tiles...), Solution: append([][2]int{}, board.Solution...)}
	return copy
}

func (r *Room) Handle(playerID string, message ClientMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	player := r.participantLocked(playerID)
	if player == nil {
		return errors.New("player is not in this mahjong room")
	}
	switch message.Type {
	case "pair":
		if err := r.removePairLocked(player, message.First, message.Second); err != nil {
			return err
		}
	case "hint":
		if err := r.requestHintLocked(player); err != nil {
			return err
		}
	case "chat":
		if err := r.addChatLocked(player, message.Content); err != nil {
			return err
		}
	case "rematch":
		if player.Seat < 0 {
			return errors.New("spectators cannot ready")
		}
		if r.status == "playing" {
			return errors.New("the round is still active")
		}
		r.ready[player.Seat] = true
		if r.ready[0] && r.ready[1] {
			if err := r.startRoundLocked(); err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported mahjong message")
	}
	r.broadcastLocked()
	return nil
}

func (r *Room) removePairLocked(player *Player, first, second int) error {
	if player.Seat < 0 {
		return errors.New("spectators cannot remove tiles")
	}
	if r.status != "playing" {
		return errors.New("the round has finished")
	}
	state := r.stateIndex(player.Seat)
	board, removed := r.boards[state], r.removed[state]
	if !board.CanRemove(removed, first, second) {
		return errors.New("select two free matching tiles")
	}
	removed[first], removed[second] = true, true
	r.pairs[player.Seat]++
	r.hint[player.Seat] = nil
	if r.mode != "duel" {
		r.removedBy[first], r.removedBy[second] = player.Seat, player.Seat
	}
	remaining := board.Remaining(removed)
	if remaining == 0 {
		r.finishClearLocked(player.Seat)
		return nil
	}
	if _, _, available := board.AvailablePair(removed); !available {
		if err := board.Reshuffle(removed); err != nil {
			return err
		}
		r.reshuffles++
		r.message = "No pairs remained, so the open tiles were reshuffled"
	} else if r.mode == "duel" {
		r.message = player.Name + " cleared a pair"
	} else {
		r.message = "The garden is opening up"
	}
	return nil
}

func (r *Room) requestHintLocked(player *Player) error {
	if player.Seat < 0 {
		return errors.New("spectators cannot request hints")
	}
	if r.status != "playing" {
		return errors.New("the round has finished")
	}
	if r.hints[player.Seat] <= 0 {
		return errors.New("no hints remaining")
	}
	state := r.stateIndex(player.Seat)
	first, second, ok := r.boards[state].AvailablePair(r.removed[state])
	if !ok {
		return errors.New("no matching pair is currently available")
	}
	r.hints[player.Seat]--
	r.hint[player.Seat] = []int{first, second}
	return nil
}

func (r *Room) finishClearLocked(seat int) {
	r.status = "finished"
	if r.timer != nil {
		r.timer.Stop()
	}
	r.deadline = nil
	if r.mode == "duel" {
		r.winner = seat
		r.message = r.players[seat].Name + " cleared the board first"
	} else {
		r.winner = -1
		r.message = "The garden was cleared together"
	}
}

func (r *Room) finishTimedLocked() {
	r.status, r.deadline = "finished", nil
	if r.mode == "garden" {
		r.winner = -1
		r.message = "Time expired before the garden was cleared"
		return
	}
	left := r.boards[0].Remaining(r.removed[0])
	right := r.boards[1].Remaining(r.removed[1])
	switch {
	case left < right:
		r.winner = 0
	case right < left:
		r.winner = 1
	default:
		r.winner = -1
	}
	if r.winner >= 0 {
		r.message = r.players[r.winner].Name + " had fewer tiles when time expired"
	} else {
		r.message = "The duel ended level"
	}
}

func (r *Room) stateIndex(seat int) int {
	if r.mode == "duel" && seat == 1 {
		return 1
	}
	return 0
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
	viewerSeat := -1
	if viewer != nil {
		viewerSeat = viewer.Seat
	}
	state := r.stateIndex(viewerSeat)
	board, removed := r.boards[state], r.removed[state]
	s := &Snapshot{RoomID: r.id, Mode: r.mode, Status: r.status, Tiles: append([]Tile{}, board.Tiles...), Removed: append([]bool{}, removed...), Free: []int{}, ViewerSeat: viewerSeat, Winner: r.winner, Message: r.message, Chat: append([]ChatMessage{}, r.chat...), Spectators: len(r.spectators), Ready: []int{}, InviteCode: r.inviteCode, StartedAt: r.startedAt, Deadline: cloneTime(r.deadline), Reshuffles: r.reshuffles}
	if r.mode != "duel" {
		s.RemovedBy = append([]int{}, r.removedBy...)
	} else {
		s.RemovedBy = []int{}
	}
	for index := range board.Tiles {
		if board.IsFree(removed, index) {
			s.Free = append(s.Free, index)
		}
	}
	for seat, player := range r.players {
		seatState := r.stateIndex(seat)
		s.Seats[seat] = &SeatView{ID: player.ID, Name: player.Name, Pairs: r.pairs[seat], Remaining: r.boards[seatState].Remaining(r.removed[seatState]), Connected: player.connected.Load()}
	}
	if viewerSeat >= 0 {
		s.HintsRemaining = r.hints[viewerSeat]
		s.Hint = append([]int{}, r.hint[viewerSeat]...)
	} else {
		s.Hint = []int{}
	}
	for seat, ready := range r.ready {
		if ready {
			s.Ready = append(s.Ready, seat)
		}
	}
	if s.Chat == nil {
		s.Chat = []ChatMessage{}
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

func (r *Room) SendState(player *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sendLocked(player)
}

func (r *Room) participantLocked(id string) *Player {
	for _, player := range r.players {
		if player != nil && player.ID == id {
			return player
		}
	}
	return r.spectators[id]
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
