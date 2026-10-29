package arena

import (
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
	key := matchKey(options)
	waiting := h.waiting[key]
	if waiting == nil || !waiting.connected.Load() || waiting.ID == player.ID {
		if waiting != nil && waiting.ID == player.ID {
			waiting.connected.Store(false)
		}
		player.Mark = "one"
		h.waiting[key] = player
		return player, nil, nil
	}

	game, gameErr := NewGameWithMode(options.Kind, options.Mode)
	if gameErr != nil {
		return nil, nil, gameErr
	}
	waiting.Mark, player.Mark = "one", "two"
	room := newRoom(uuid.NewString(), game, waiting, player, settingsFrom(options), h.scheduleCleanup)
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
	if player != nil && existing.status != "archived" {
		player.replaceSender()
		existing.updated = time.Now()
		return player, existing
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

func matchKey(options JoinOptions) MatchKey {
	return MatchKey{Kind: options.Kind, Mode: options.Mode, BestOf: options.BestOf, TimeControl: options.TimeControl, InviteCode: options.InviteCode}
}

func settingsFrom(options JoinOptions) RoomSettings {
	visibility := "public"
	if options.InviteCode != "" {
		visibility = "private"
	}
	return RoomSettings{Mode: options.Mode, BestOf: options.BestOf, TimeControl: options.TimeControl, Visibility: visibility, InviteCode: options.InviteCode}
}
