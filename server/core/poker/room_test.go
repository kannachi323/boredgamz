package poker

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnapshotCollectionsSerializeAsArrays(t *testing.T) {
	alice, bob := newPlayer("json-a", "Alice"), newPlayer("json-b", "Bob")
	room, err := newRoom("json-room", "holdem", alice, bob, "")
	require.NoError(t, err)
	room.mu.Lock()
	payload, err := json.Marshal(room.snapshotLocked(bob))
	room.mu.Unlock()
	require.NoError(t, err)
	require.NotContains(t, string(payload), `"community":null`)
	require.NotContains(t, string(payload), `"legalActions":null`)
	require.NotContains(t, string(payload), `"chat":null`)
}

func TestSnapshotsNeverLeakOpponentHoleCards(t *testing.T) {
	alice, bob := newPlayer("alice", "Alice"), newPlayer("bob", "Bob")
	room, err := newRoom("private-cards", "holdem", alice, bob, "")
	require.NoError(t, err)
	room.mu.Lock()
	aliceView := room.snapshotLocked(alice)
	bobView := room.snapshotLocked(bob)
	spectatorView := room.snapshotLocked(nil)
	room.mu.Unlock()
	require.Len(t, aliceView.Seats[0].Cards, 2)
	require.Empty(t, aliceView.Seats[1].Cards)
	require.Len(t, bobView.Seats[1].Cards, 2)
	require.Empty(t, bobView.Seats[0].Cards)
	require.Empty(t, spectatorView.Seats[0].Cards)
	require.Empty(t, spectatorView.Seats[1].Cards)
	// Clients find their seat from viewerSeat because authenticated seat IDs
	// are rewritten server-side and never equal the client's local ID.
	require.Equal(t, 0, aliceView.ViewerSeat)
	require.Equal(t, "alice", aliceView.ViewerID)
	require.Equal(t, 1, bobView.ViewerSeat)
	require.Equal(t, -1, spectatorView.ViewerSeat)
}

func TestShowdownRevealsBothHands(t *testing.T) {
	alice, bob := newPlayer("show-a", "Alice"), newPlayer("show-b", "Bob")
	room, err := newRoom("showdown", "holdem", alice, bob, "")
	require.NoError(t, err)
	room.mu.Lock()
	room.table.Phase = "showdown"
	view := room.snapshotLocked(nil)
	room.mu.Unlock()
	// Once a showdown is declared, both contested hands are public.
	require.Len(t, view.Seats[0].Cards, 2)
	require.Len(t, view.Seats[1].Cards, 2)
	room.mu.Lock()
	playerView := room.snapshotLocked(alice)
	room.mu.Unlock()
	require.Len(t, playerView.Seats[0].Cards, 2)
	require.Len(t, playerView.Seats[1].Cards, 2)
}
