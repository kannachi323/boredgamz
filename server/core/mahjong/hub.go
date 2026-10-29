package mahjong

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type matchKey struct{ mode, invite string }

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
		if player != nil {
			player.replaceSender()
			room.mu.Unlock()
			return player, room, nil
		}
		room.mu.Unlock()
		delete(h.playerRooms, input.PlayerID)
	}
	player := newPlayer(input.PlayerID, input.Name)
	key := matchKey{input.Mode, input.InviteCode}
	waiting := h.waiting[key]
	if waiting == nil || !waiting.connected.Load() || waiting.ID == player.ID {
		player.Seat = 0
		h.waiting[key] = player
		return player, nil, nil
	}
	room, err := newRoom(uuid.NewString(), input.Mode, waiting, player, input.InviteCode)
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
	player.connected.Store(false)
	if player.Seat < 0 {
		delete(room.spectators, player.ID)
	} else if room.mode == "duel" && room.status == "playing" {
		room.status, room.winner = "finished", 1-player.Seat
		room.message = room.players[1-player.Seat].Name + " wins after the opponent left"
		if room.timer != nil {
			room.timer.Stop()
		}
		room.deadline = nil
	}
	room.broadcastLocked()
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
