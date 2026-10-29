package arena

import "time"

func (h *Hub) Disconnect(player *Player, room *Room, generation uint64) {
	if player == nil || room == nil || player.generation.Load() != generation {
		return
	}

	room.mu.Lock()
	if player.generation.Load() != generation {
		room.mu.Unlock()
		return
	}
	player.connected.Store(false)
	if player.Mark == "spectator" {
		delete(room.spectators, player.ID)
		room.broadcastStateLocked()
		room.mu.Unlock()
		h.mu.Lock()
		if h.playerRooms[player.ID] == room {
			delete(h.playerRooms, player.ID)
		}
		h.mu.Unlock()
		return
	}
	room.broadcastStateLocked()
	room.mu.Unlock()
	h.scheduleAbandonCheck(room)

	time.AfterFunc(disconnectGrace, func() {
		room.mu.Lock()
		defer room.mu.Unlock()
		if player.generation.Load() != generation || player.connected.Load() || room.status != "playing" {
			return
		}
		room.finishRoundLocked(opposite(player.Mark))
		room.updated = time.Now()
		room.broadcastStateLocked()
	})
}
