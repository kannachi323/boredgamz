package poker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPokerMatchmakingSeparatesModesAndPrivateCodes(t *testing.T) {
	hub := NewHub()
	holdem, room, err := hub.Join(JoinOptions{PlayerID: "holdem-one", Name: "A", Mode: "holdem", InviteCode: "river"})
	require.NoError(t, err)
	require.Nil(t, room)
	_, drawRoom, err := hub.Join(JoinOptions{PlayerID: "draw-one", Name: "B", Mode: "five-card-draw", InviteCode: "river"})
	require.NoError(t, err)
	require.Nil(t, drawRoom)
	_, wrongCode, err := hub.Join(JoinOptions{PlayerID: "holdem-wrong", Name: "C", Mode: "holdem", InviteCode: "other"})
	require.NoError(t, err)
	require.Nil(t, wrongCode)
	second, matched, err := hub.Join(JoinOptions{PlayerID: "holdem-two", Name: "D", Mode: "holdem", InviteCode: "RIVER"})
	require.NoError(t, err)
	require.NotNil(t, matched)
	require.Equal(t, "holdem", matched.table.Mode)
	require.Equal(t, holdem.ID, matched.players[0].ID)
	require.Equal(t, second.ID, matched.players[1].ID)
}

func TestPokerReconnectKeepsSeatAndPrivateCards(t *testing.T) {
	hub := NewHub()
	first, _, err := hub.Join(JoinOptions{PlayerID: "reconnect-a", Name: "Alice", Mode: "holdem"})
	require.NoError(t, err)
	_, room, err := hub.Join(JoinOptions{PlayerID: "reconnect-b", Name: "Bob", Mode: "holdem"})
	require.NoError(t, err)
	generation := first.generation.Load()
	first.connected.Store(false)
	rejoined, sameRoom, err := hub.Join(JoinOptions{PlayerID: "reconnect-a", Name: "Alice", Mode: "holdem"})
	require.NoError(t, err)
	require.Same(t, room, sameRoom)
	require.Equal(t, 0, rejoined.Seat)
	require.Greater(t, rejoined.generation.Load(), generation)
	room.mu.Lock()
	snapshot := room.snapshotLocked(rejoined)
	room.mu.Unlock()
	require.NotEmpty(t, snapshot.Seats[0].Cards)
	require.Empty(t, snapshot.Seats[1].Cards)
}

func TestPokerSpectatorReceivesNoActiveHoleCards(t *testing.T) {
	hub := NewHub()
	_, _, err := hub.Join(JoinOptions{PlayerID: "spectate-a", Name: "Alice", Mode: "holdem"})
	require.NoError(t, err)
	_, room, err := hub.Join(JoinOptions{PlayerID: "spectate-b", Name: "Bob", Mode: "holdem"})
	require.NoError(t, err)
	watcher, watchedRoom, err := hub.Join(JoinOptions{PlayerID: "watcher", Name: "Rail", Role: "spectator", RoomID: room.id})
	require.NoError(t, err)
	require.Same(t, room, watchedRoom)
	require.Equal(t, -1, watcher.Seat)
	room.mu.Lock()
	snapshot := room.snapshotLocked(watcher)
	room.mu.Unlock()
	require.Empty(t, snapshot.Seats[0].Cards)
	require.Empty(t, snapshot.Seats[1].Cards)
}

func TestPokerEmptyRoomCleanupReleasesAllMappings(t *testing.T) {
	hub := NewHub()
	first, _, err := hub.Join(JoinOptions{PlayerID: "cleanup-a", Name: "Alice", Mode: "holdem"})
	require.NoError(t, err)
	second, room, err := hub.Join(JoinOptions{PlayerID: "cleanup-b", Name: "Bob", Mode: "holdem"})
	require.NoError(t, err)
	first.connected.Store(false)
	second.connected.Store(false)
	hub.cleanupRoomIfEmpty(room)
	require.NotContains(t, hub.rooms, room.id)
	require.NotContains(t, hub.playerRooms, first.ID)
	require.NotContains(t, hub.playerRooms, second.ID)
}
