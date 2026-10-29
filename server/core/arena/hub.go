package arena

import (
	"boredgamz/core/identity"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	mu          sync.Mutex
	waiting     map[MatchKey]*Player
	rooms       map[string]*Room
	playerRooms map[string]*Room
}

type HubStats struct {
	Waiting     int `json:"waiting"`
	Rooms       int `json:"rooms"`
	Connections int `json:"connections"`
}

func NewHub() *Hub {
	return &Hub{
		waiting: make(map[MatchKey]*Player), rooms: make(map[string]*Room),
		playerRooms: make(map[string]*Room),
	}
}

func (h *Hub) Stats() HubStats {
	h.mu.Lock()
	defer h.mu.Unlock()
	return HubStats{Waiting: len(h.waiting), Rooms: len(h.rooms), Connections: len(h.playerRooms)}
}

func ValidKind(kind GameKind) bool {
	return kind == TicTacToe || kind == ConnectFour || kind == Gomoku || kind == Chess
}

func SanitizeName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Guest"
	}
	runes := []rune(name)
	if len(runes) > 24 {
		runes = runes[:24]
	}
	return string(runes)
}

func (h *Hub) Join(kind GameKind, playerID, playerName string) (*Player, *Room, error) {
	return h.JoinMatch(JoinOptions{Kind: kind, PlayerID: playerID, PlayerName: playerName})
}

func (h *Hub) JoinMatch(options JoinOptions) (*Player, *Room, error) {
	options, err := normalizeJoinOptions(options)
	if err != nil {
		return nil, nil, err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if options.Role == "spectator" {
		return h.joinSpectatorLocked(options)
	}
	if player, room := h.reconnectLocked(options); player != nil {
		return player, room, nil
	}

	player := newPlayer(options.PlayerID, options.PlayerName)
	player.options = options
	key := matchKey(options)
	waiting := h.waiting[key]
	if waiting == nil || !waiting.connected.Load() || identity.SameAccount(waiting.ID, player.ID) {
		if waiting != nil && identity.SameAccount(waiting.ID, player.ID) {
			waiting.connected.Store(false)
		}
		player.Mark = "one"
		h.waiting[key] = player
		return player, nil, nil
	}

	host := waiting.options
	game, gameErr := NewGameWithMode(host.Kind, host.Mode)
	if gameErr != nil {
		return nil, nil, gameErr
	}
	waiting.Mark, player.Mark = "one", "two"
	room := newRoom(uuid.NewString(), game, waiting, player, settingsFrom(host), h.scheduleCleanup)
	delete(h.waiting, key)
	h.rooms[room.id] = room
	h.playerRooms[waiting.ID], h.playerRooms[player.ID] = room, room
	return player, room, nil
}

func (h *Hub) reconnectLocked(options JoinOptions) (*Player, *Room) {
	existing := h.playerRooms[options.PlayerID]
	if existing == nil {
		return nil, nil
	}
	if existing.game.Kind != options.Kind {
		delete(h.playerRooms, options.PlayerID)
		return nil, nil
	}
	existing.mu.Lock()
	defer existing.mu.Unlock()
	player := existing.participantByIDLocked(options.PlayerID)
	// A game in progress always resumes (you can't dodge a loss by asking for
	// another table). Between rounds, a request for a different table means
	// the player moved on: release the old seat and queue for the new one.
	sameTable := matchKey(options) == matchKey(JoinOptions{Kind: existing.game.Kind, Mode: existing.settings.Mode, BestOf: existing.settings.BestOf, TimeControl: existing.settings.TimeControl, InviteCode: existing.settings.InviteCode})
	if player != nil && existing.status != "archived" && (existing.status == "playing" || sameTable) {
		player.replaceSender()
		existing.updated = time.Now()
		return player, existing
	}
	if player != nil && player.Mark != "spectator" {
		player.connected.Store(false)
		player.left.Store(true)
		existing.updated = time.Now()
		existing.broadcastStateLocked()
	}
	delete(h.playerRooms, options.PlayerID)
	return nil, nil
}

func (h *Hub) LeaveRoom(player *Player, room *Room) {
	if player == nil || room == nil {
		return
	}
	room.mu.Lock()
	player.connected.Store(false)
	player.left.Store(true)
	if player.Mark == "spectator" {
		delete(room.spectators, player.ID)
	} else if room.status == "playing" {
		room.finishRoundLocked(opposite(player.Mark))
	}
	room.updated = time.Now()
	room.broadcastStateLocked()
	room.mu.Unlock()

	h.mu.Lock()
	if h.playerRooms[player.ID] == room {
		delete(h.playerRooms, player.ID)
	}
	h.mu.Unlock()
}

func (h *Hub) joinSpectatorLocked(options JoinOptions) (*Player, *Room, error) {
	room := h.rooms[options.RoomID]
	if room == nil || room.game.Kind != options.Kind {
		return nil, nil, errors.New("room not found")
	}
	player := newPlayer(options.PlayerID, options.PlayerName)
	player.Mark = "spectator"
	room.mu.Lock()
	room.spectators[player.ID] = player
	room.updated = time.Now()
	room.mu.Unlock()
	h.playerRooms[player.ID] = room
	return player, room, nil
}

func (h *Hub) LeaveWaiting(options JoinOptions, player *Player) {
	h.mu.Lock()
	defer h.mu.Unlock()
	key := matchKey(options)
	if h.waiting[key] == player {
		delete(h.waiting, key)
	}
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

func (h *Hub) scheduleCleanup(room *Room) {
	room.cleanupOnce.Do(func() {
		time.AfterFunc(roomRetention, func() {
			h.mu.Lock()
			defer h.mu.Unlock()
			delete(h.rooms, room.id)
			for _, player := range room.allParticipants() {
				if h.playerRooms[player.ID] == room {
					delete(h.playerRooms, player.ID)
				}
			}
		})
	})
}

func normalizeJoinOptions(options JoinOptions) (JoinOptions, error) {
	if !ValidKind(options.Kind) {
		return options, errors.New("unsupported game")
	}
	if options.PlayerID == "" || len(options.PlayerID) > 80 {
		return options, errors.New("invalid player ID")
	}
	options.Mode = NormalizeMode(options.Kind, options.Mode)
	if options.BestOf != 3 && options.BestOf != 5 && options.BestOf != 7 {
		options.BestOf = 1
	}
	if options.Kind == Chess {
		switch options.TimeControl {
		case "untimed", "3m+2s", "10m", "15m+10s":
		default:
			options.TimeControl = "10m"
		}
	} else {
		switch options.TimeControl {
		case "15s", "30s", "60s":
		default:
			options.TimeControl = "untimed"
		}
	}
	options.InviteCode = strings.ToUpper(strings.TrimSpace(options.InviteCode))
	if len(options.InviteCode) > 12 {
		return options, errors.New("invite code is too long")
	}
	if options.Role == "spectator" && options.RoomID == "" {
		return options, errors.New("spectators need a room ID")
	}
	return options, nil
}

// matchKey queues public players by exact settings. A private code alone
// identifies the table: whoever joins with it plays the host's settings.
func matchKey(options JoinOptions) MatchKey {
	if options.InviteCode != "" {
		return MatchKey{Kind: options.Kind, InviteCode: options.InviteCode}
	}
	return MatchKey{Kind: options.Kind, Mode: options.Mode, BestOf: options.BestOf, TimeControl: options.TimeControl, InviteCode: options.InviteCode}
}

func settingsFrom(options JoinOptions) RoomSettings {
	visibility := "public"
	if options.InviteCode != "" {
		visibility = "private"
	}
	return RoomSettings{Mode: options.Mode, BestOf: options.BestOf, TimeControl: options.TimeControl, Visibility: visibility, InviteCode: options.InviteCode}
}

// assignedRoom returns the room matchmaking put this player in, if any. A
// socket can close before its loop notices the match; it must still count
// as a disconnect from that room, not as leaving the queue.
func (h *Hub) assignedRoom(player *Player) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.playerRooms[player.ID]
}
