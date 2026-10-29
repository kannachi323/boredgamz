package mahjong

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHubSeparatesModesAndPrivateCodes(t *testing.T) {
	hub := NewHub()
	first, room, err := hub.Join(JoinOptions{PlayerID: "a", Name: "Alice", Mode: "garden", InviteCode: "LEAF"})
	require.NoError(t, err)
	require.Nil(t, room)
	_, room, err = hub.Join(JoinOptions{PlayerID: "b", Name: "Bob", Mode: "duel", InviteCode: "LEAF"})
	require.NoError(t, err)
	require.Nil(t, room)
	second, room, err := hub.Join(JoinOptions{PlayerID: "c", Name: "Cora", Mode: "garden", InviteCode: "leaf"})
	require.NoError(t, err)
	require.NotNil(t, room)
	require.Equal(t, first, room.players[0])
	require.Equal(t, second, room.players[1])
	room.timer.Stop()
}

func TestReconnectReturnsOriginalSeat(t *testing.T) {
	hub := NewHub()
	_, _, err := hub.Join(JoinOptions{PlayerID: "a", Name: "Alice", Mode: "zen"})
	require.NoError(t, err)
	_, room, err := hub.Join(JoinOptions{PlayerID: "b", Name: "Bob", Mode: "zen"})
	require.NoError(t, err)
	original := room.players[0]
	original.connected.Store(false)
	reconnected, sameRoom, err := hub.Join(JoinOptions{PlayerID: "a", Name: "Alice", Mode: "zen"})
	require.NoError(t, err)
	require.Same(t, original, reconnected)
	require.Same(t, room, sameRoom)
	require.True(t, reconnected.connected.Load())
}
