package mahjong

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func testRoom(t *testing.T, mode string) *Room {
	t.Helper()
	room, err := newRoom("room", mode, newPlayer("a", "Alice"), newPlayer("b", "Bob"), "")
	require.NoError(t, err)
	t.Cleanup(func() {
		room.mu.Lock()
		defer room.mu.Unlock()
		if room.timer != nil {
			room.timer.Stop()
		}
	})
	return room
}

func TestGardenPairIsSharedAndAttributed(t *testing.T) {
	room := testRoom(t, "garden")
	pair := room.boards[0].Solution[0]
	require.NoError(t, room.Handle("a", ClientMessage{Type: "pair", First: pair[0], Second: pair[1]}))
	room.mu.Lock()
	defer room.mu.Unlock()
	require.True(t, room.removed[0][pair[0]])
	require.True(t, room.removed[1][pair[1]])
	require.Equal(t, 0, room.removedBy[pair[0]])
	require.Equal(t, 1, room.pairs[0])
	require.Equal(t, 142, room.snapshotLocked(room.players[1]).Seats[1].Remaining)
}

func TestDuelBoardsProgressIndependently(t *testing.T) {
	room := testRoom(t, "duel")
	pair := room.boards[0].Solution[0]
	require.NoError(t, room.Handle("a", ClientMessage{Type: "pair", First: pair[0], Second: pair[1]}))
	room.mu.Lock()
	defer room.mu.Unlock()
	require.True(t, room.removed[0][pair[0]])
	require.False(t, room.removed[1][pair[0]])
	require.Equal(t, 70, room.snapshotLocked(room.players[0]).Seats[0].Remaining)
	require.Equal(t, 72, room.snapshotLocked(room.players[0]).Seats[1].Remaining)
}

func TestSimultaneousGardenPairCanOnlyBeClaimedOnce(t *testing.T) {
	room := testRoom(t, "garden")
	pair := room.boards[0].Solution[0]
	errors := make(chan error, 2)
	var wait sync.WaitGroup
	for _, id := range []string{"a", "b"} {
		wait.Add(1)
		go func(playerID string) {
			defer wait.Done()
			errors <- room.Handle(playerID, ClientMessage{Type: "pair", First: pair[0], Second: pair[1]})
		}(id)
	}
	wait.Wait()
	close(errors)
	succeeded := 0
	for err := range errors {
		if err == nil {
			succeeded++
		}
	}
	require.Equal(t, 1, succeeded)
	require.Equal(t, 1, room.pairs[0]+room.pairs[1])
}

func TestHintIsPrivateFreeAndLimited(t *testing.T) {
	room := testRoom(t, "zen")
	require.NoError(t, room.Handle("a", ClientMessage{Type: "hint"}))
	room.mu.Lock()
	defer room.mu.Unlock()
	alice := room.snapshotLocked(room.players[0])
	bob := room.snapshotLocked(room.players[1])
	require.Len(t, alice.Hint, 2)
	require.True(t, room.boards[0].CanRemove(room.removed[0], alice.Hint[0], alice.Hint[1]))
	require.Empty(t, bob.Hint)
	require.Equal(t, 2, alice.HintsRemaining)
	require.Equal(t, 3, bob.HintsRemaining)
}

func TestDuelTimerChoosesBoardWithFewerRemainingTiles(t *testing.T) {
	room := testRoom(t, "duel")
	pair := room.boards[0].Solution[0]
	require.NoError(t, room.Handle("a", ClientMessage{Type: "pair", First: pair[0], Second: pair[1]}))
	room.mu.Lock()
	room.finishTimedLocked()
	require.Equal(t, "finished", room.status)
	require.Equal(t, 0, room.winner)
	room.mu.Unlock()
}

func TestBothPlayersMustReadyBeforeFreshBoard(t *testing.T) {
	room := testRoom(t, "duel")
	room.mu.Lock()
	room.status = "finished"
	original := room.boards[0]
	room.mu.Unlock()
	require.NoError(t, room.Handle("a", ClientMessage{Type: "rematch"}))
	require.Same(t, original, room.boards[0])
	require.NoError(t, room.Handle("b", ClientMessage{Type: "rematch"}))
	require.NotSame(t, original, room.boards[0])
	require.Equal(t, "playing", room.status)
}
