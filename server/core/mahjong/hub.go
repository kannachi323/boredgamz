package mahjong

import (
	"boredgamz/core/identity"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type matchKey struct{ mode, invite string }

// keyFor queues public players by mode. A private code alone identifies the
// table: whoever joins with it plays the host's mode.
func keyFor(mode, invite string) matchKey {
	if invite != "" {
		return matchKey{invite: invite}
	}
	return matchKey{mode: mode}
}

type Hub struct {
	mu          sync.Mutex
	waiting     map[matchKey]*Player
	rooms       map[string]*Room
	playerRooms map[string]*Room
}

type HubStats struct {
	Waiting     int `json:"waiting"`
	Rooms       int `json:"rooms"`
	Connections int `json:"connections"`
}

func NewHub() *Hub {
	return &Hub{waiting: map[matchKey]*Player{}, rooms: map[string]*Room{}, playerRooms: map[string]*Room{}}
}

func (h *Hub) Stats() HubStats {
	h.mu.Lock()
	defer h.mu.Unlock()
	return HubStats{Waiting: len(h.waiting), Rooms: len(h.rooms), Connections: len(h.playerRooms)}
}

func (h *Hub) Join(input JoinOptions) (*Player, *Room, error) {
	input.Mode = normalizeMode(input.Mode)
	input.InviteCode = strings.ToUpper(strings.TrimSpace(input.InviteCode))
	if len(input.InviteCode) > 12 {
		return nil, nil, errors.New("invite code is too long")
	}
	if input.PlayerID == "" || len(input.PlayerID) > 80 {
		return nil, nil, errors.New("invalid player ID")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if input.Role == "spectator" {
		room := h.rooms[input.RoomID]
		if room == nil {
			return nil, nil, errors.New("mahjong room not found")
		}
		player := newPlayer(input.PlayerID, input.Name)
		room.mu.Lock()
		room.spectators[player.ID] = player
		room.mu.Unlock()
		h.playerRooms[player.ID] = room
		return player, room, nil
	}
	if room := h.playerRooms[input.PlayerID]; room != nil {
		room.mu.Lock()
		player := room.participantLocked(input.PlayerID)
		// Same table: this is a reconnect. A different table means the player
		// moved on, so they leave the old one first.
		if player != nil && keyFor(room.mode, room.inviteCode) == keyFor(normalizeMode(input.Mode), input.InviteCode) {
			player.replaceSender()
			room.mu.Unlock()
			return player, room, nil
		}
		if player != nil {
			room.leaveLocked(player)
		}
		room.mu.Unlock()
		delete(h.playerRooms, input.PlayerID)
		if player != nil && player.Seat >= 0 {
			h.scheduleCleanup(room)
		}
	}
	player := newPlayer(input.PlayerID, input.Name)
	player.mode = input.Mode
	key := keyFor(input.Mode, input.InviteCode)
	waiting := h.waiting[key]
	if waiting == nil || !waiting.connected.Load() || identity.SameAccount(waiting.ID, player.ID) {
		player.Seat = 0
		h.waiting[key] = player
		return player, nil, nil
	}
	room, err := newRoom(uuid.NewString(), waiting.mode, waiting, player, input.InviteCode)
	if err != nil {
		return nil, nil, err
	}
	delete(h.waiting, key)
	h.rooms[room.id] = room
	h.playerRooms[waiting.ID], h.playerRooms[player.ID] = room, room
	return player, room, nil
}

func (h *Hub) MatchWaiting(player *Player) *Room {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		h.mu.Lock()
		room := h.playerRooms[player.ID]
		h.mu.Unlock()
		if room != nil {
			return room
		}
		if !player.connected.Load() {
			return nil
		}
	}
	return nil
}

func (h *Hub) LeaveWaiting(player *Player) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for key, waiting := range h.waiting {
		if waiting == player {
			delete(h.waiting, key)
		}
	}
}

func (h *Hub) Leave(player *Player, room *Room) {
	if player == nil || room == nil {
		return
	}
	room.mu.Lock()
	room.leaveLocked(player)
	room.mu.Unlock()
	h.mu.Lock()
	if h.playerRooms[player.ID] == room {
		delete(h.playerRooms, player.ID)
	}
	h.mu.Unlock()
	if player.Seat >= 0 {
		h.scheduleCleanup(room)
	}
}

func (h *Hub) Disconnect(player *Player, room *Room, generation uint64) {
	if player == nil || room == nil || player.generation.Load() != generation {
		return
	}
	room.mu.Lock()
	player.connected.Store(false)
	if player.Seat < 0 {
		delete(room.spectators, player.ID)
	}
	room.broadcastLocked()
	room.mu.Unlock()
	if player.Seat < 0 {
		h.mu.Lock()
		if h.playerRooms[player.ID] == room {
			delete(h.playerRooms, player.ID)
		}
		h.mu.Unlock()
		return
	}
	h.scheduleCleanup(room)
	time.AfterFunc(20*time.Second, func() {
		room.mu.Lock()
		defer room.mu.Unlock()
		if player.generation.Load() != generation || player.connected.Load() || room.status != "playing" || room.mode != "duel" {
			return
		}
		room.status, room.winner = "finished", 1-player.Seat
		room.message = room.players[1-player.Seat].Name + " wins after the opponent disconnected"
		if room.timer != nil {
			room.timer.Stop()
		}
		room.deadline = nil
		room.broadcastLocked()
	})
}

func (h *Hub) scheduleCleanup(room *Room) {
	time.AfterFunc(5*time.Minute, func() { h.cleanupRoomIfEmpty(room) })
}

func (h *Hub) cleanupRoomIfEmpty(room *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.players[0].connected.Load() || room.players[1].connected.Load() {
		return
	}
	if room.timer != nil {
		room.timer.Stop()
	}
	delete(h.rooms, room.id)
	for _, player := range room.players {
		if h.playerRooms[player.ID] == room {
			delete(h.playerRooms, player.ID)
		}
	}
}

// assignedRoom returns the room matchmaking put this player in, if any. A
// socket can close before its loop notices the match; it must still count
// as a disconnect from that room, not as leaving the queue.
func (h *Hub) assignedRoom(player *Player) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.playerRooms[player.ID]
}

// leaveLocked removes a participant on purpose: leaving a live duel concedes
// it. The caller holds r.mu.
func (r *Room) leaveLocked(player *Player) {
	player.connected.Store(false)
	player.left.Store(true)
	if player.Seat < 0 {
		delete(r.spectators, player.ID)
	} else if r.mode == "duel" && r.status == "playing" {
		r.status, r.winner = "finished", 1-player.Seat
		r.message = r.players[1-player.Seat].Name + " wins after the opponent left"
		if r.timer != nil {
			r.timer.Stop()
		}
		r.deadline = nil
	}
	r.broadcastLocked()
}
