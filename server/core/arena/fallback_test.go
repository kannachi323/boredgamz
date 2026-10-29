package arena

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func withFallback(t *testing.T, after time.Duration) {
	t.Helper()
	previous := queueFallbackAfter
	queueFallbackAfter = after
	t.Cleanup(func() { queueFallbackAfter = previous })
}

// matchWithin runs the queue poller for a waiting player and returns its room.
func matchWithin(hub *Hub, player *Player, limit time.Duration) *Room {
	done := make(chan *Room, 1)
	go func() { done <- hub.MatchWaiting(player) }()
	select {
	case room := <-done:
		return room
	case <-time.After(limit):
		player.connected.Store(false) // stop the poller
		<-done
		return nil
	}
}

func TestFallbackPairsDifferentSettingsInTheSameMode(t *testing.T) {
	withFallback(t, 150*time.Millisecond)
	hub := NewHub()
	early, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, Mode: "classic", BestOf: 3, TimeControl: "30s", PlayerID: "early"})
	require.NoError(t, err)
	require.Nil(t, room)
	time.Sleep(20 * time.Millisecond)
	late, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, Mode: "classic", BestOf: 5, TimeControl: "untimed", PlayerID: "late"})
	require.NoError(t, err)
	require.Nil(t, room, "different settings never match immediately")

	matched := matchWithin(hub, late, time.Second)
	require.NotNil(t, matched)
	require.Same(t, matched, hub.playerRooms[early.ID])
	snapshot := matched.Snapshot()
	require.Equal(t, 3, snapshot.Settings.BestOf, "the player who waited longer hosts")
	require.Equal(t, "30s", snapshot.Settings.TimeControl)
	require.Equal(t, "early", snapshot.Players[0].ID)
	require.Zero(t, hub.Stats().Waiting)
	for _, player := range []*Player{early, late} {
		select {
		case message := <-player.sender():
			require.Equal(t, "state", message.Type, player.ID)
		default:
			t.Fatalf("%s was not sent the new table", player.ID)
		}
	}
}

func TestFallbackWaitsForTheDelay(t *testing.T) {
	withFallback(t, time.Hour)
	hub := NewHub()
	_, _, err := hub.JoinMatch(JoinOptions{Kind: ConnectFour, BestOf: 3, PlayerID: "x"})
	require.NoError(t, err)
	late, _, err := hub.JoinMatch(JoinOptions{Kind: ConnectFour, BestOf: 5, PlayerID: "y"})
	require.NoError(t, err)
	require.Nil(t, matchWithin(hub, late, 300*time.Millisecond))
}

func TestFallbackNeverCrossesModesPrivateTablesOrAccounts(t *testing.T) {
	withFallback(t, 50*time.Millisecond)
	hub := NewHub()
	_, _, err := hub.JoinMatch(JoinOptions{Kind: Gomoku, Mode: "freestyle", PlayerID: "free"})
	require.NoError(t, err)
	_, _, err = hub.JoinMatch(JoinOptions{Kind: Gomoku, Mode: "compact", PlayerID: "compact"})
	require.NoError(t, err)
	_, _, err = hub.JoinMatch(JoinOptions{Kind: Gomoku, Mode: "freestyle", InviteCode: "312345", PlayerID: "private"})
	require.NoError(t, err)
	_, _, err = hub.JoinMatch(JoinOptions{Kind: Chess, Mode: "rapid", PlayerID: "acct:tab1"})
	require.NoError(t, err)
	time.Sleep(60 * time.Millisecond)
	twin, _, err := hub.JoinMatch(JoinOptions{Kind: Chess, Mode: "rapid", PlayerID: "acct:tab2"})
	require.NoError(t, err)

	for _, id := range []string{"free", "compact", "private"} {
		var player *Player
		for _, waiting := range hub.waiting {
			if waiting.ID == id {
				player = waiting
			}
		}
		require.NotNil(t, player, id)
		require.Nil(t, matchWithin(hub, player, 250*time.Millisecond), id)
	}
	require.Nil(t, matchWithin(hub, twin, 250*time.Millisecond), "same account in another tab")
}
