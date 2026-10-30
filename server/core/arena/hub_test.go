package arena

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHubMatchesDistinctPlayersAndEnforcesTurns(t *testing.T) {
	hub := NewHub()
	one, firstRoom, err := hub.Join(TicTacToe, "player-one", " Alice ")
	require.NoError(t, err)
	require.Nil(t, firstRoom)
	require.Equal(t, "Alice", one.Name)

	two, room, err := hub.Join(TicTacToe, "player-two", "Bob")
	require.NoError(t, err)
	require.NotNil(t, room)
	require.Equal(t, "two", two.Mark)
	require.Equal(t, room, hub.playerRooms[one.ID])

	require.EqualError(t, room.HandleMove(two.ID, 0), "it is not your turn")
	require.NoError(t, room.HandleMove(one.ID, 0))
	require.EqualError(t, room.HandleMove(two.ID, 0), "position is already occupied")
	require.NoError(t, room.HandleMove(two.ID, 4))
	require.Equal(t, "one", room.Snapshot().Turn)
	require.NotNil(t, room.Snapshot().Chat)
	require.NotNil(t, room.Snapshot().RematchReady)
}

func TestRoomFinishesWithAuthoritativeWinner(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.Join(TicTacToe, "one", "One")
	require.NoError(t, err)
	two, room, err := hub.Join(TicTacToe, "two", "Two")
	require.NoError(t, err)

	for _, move := range []struct {
		id       string
		position int
	}{{one.ID, 0}, {two.ID, 3}, {one.ID, 1}, {two.ID, 4}, {one.ID, 2}} {
		require.NoError(t, room.HandleMove(move.id, move.position))
	}
	snapshot := room.Snapshot()
	require.Equal(t, "finished", snapshot.Status)
	require.Equal(t, "one", snapshot.Winner)
	require.EqualError(t, room.HandleMove(two.ID, 5), "game is not active")
}

func TestSameIdentityCannotMatchItself(t *testing.T) {
	hub := NewHub()
	first, room, err := hub.Join(Gomoku, "same-id", "Guest")
	require.NoError(t, err)
	require.Nil(t, room)

	second, room, err := hub.Join(Gomoku, "same-id", "Guest")
	require.NoError(t, err)
	require.Nil(t, room)
	require.NotSame(t, first, second)
	require.Equal(t, second, hub.waiting[MatchKey{Kind: Gomoku, Mode: "freestyle", BestOf: 1, TimeControl: "untimed"}])
}

func TestSameAccountInAnotherTabCannotMatchItself(t *testing.T) {
	hub := NewHub()
	_, room, err := hub.Join(TicTacToe, "account-1:tabA", "Ada")
	require.NoError(t, err)
	require.Nil(t, room)
	_, room, err = hub.Join(TicTacToe, "account-1:tabB", "Ada")
	require.NoError(t, err)
	require.Nil(t, room)
	_, room, err = hub.Join(TicTacToe, "account-2:tabA", "Bo")
	require.NoError(t, err)
	require.NotNil(t, room)
}

func TestOldConnectionCannotDisconnectReconnectedPlayer(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.Join(TicTacToe, "one", "One")
	require.NoError(t, err)
	_, room, err := hub.Join(TicTacToe, "two", "Two")
	require.NoError(t, err)
	oldGeneration := one.generation.Load()

	reconnected, sameRoom, err := hub.Join(TicTacToe, "one", "One")
	require.NoError(t, err)
	require.Same(t, one, reconnected)
	require.Same(t, room, sameRoom)
	require.Greater(t, one.generation.Load(), oldGeneration)

	hub.Disconnect(one, room, oldGeneration)
	require.True(t, room.Snapshot().Players[0].Connected)
}

func TestSanitizeNameBoundsUserInput(t *testing.T) {
	require.Equal(t, "Guest", SanitizeName("  "))
	require.Equal(t, "abcdefghijklmnopqrstuvwx", SanitizeName("abcdefghijklmnopqrstuvwxyz"))
}

func TestQueuesAreIsolatedByModeAndInviteCode(t *testing.T) {
	hub := NewHub()
	one, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, Mode: "classic", InviteCode: "table-a", PlayerID: "one"})
	require.NoError(t, err)
	require.Nil(t, room)

	// A public player with the same mode never lands on a private table.
	_, room, err = hub.JoinMatch(JoinOptions{Kind: TicTacToe, Mode: "classic", PlayerID: "public"})
	require.NoError(t, err)
	require.Nil(t, room)
	// A different code waits separately.
	_, room, err = hub.JoinMatch(JoinOptions{Kind: TicTacToe, Mode: "classic", InviteCode: "table-b", PlayerID: "other"})
	require.NoError(t, err)
	require.Nil(t, room)
	require.NotNil(t, hub.waiting[MatchKey{Kind: TicTacToe, InviteCode: "TABLE-A"}])
	require.True(t, one.connected.Load())

	// The guest only needs the code; the host's mode and series length win.
	_, room, err = hub.JoinMatch(JoinOptions{Kind: TicTacToe, Mode: "misere", BestOf: 7, InviteCode: "TABLE-A", PlayerID: "three"})
	require.NoError(t, err)
	require.NotNil(t, room)
	snapshot := room.Snapshot()
	require.Equal(t, "private", snapshot.Settings.Visibility)
	require.Equal(t, "classic", snapshot.Settings.Mode)
	require.Equal(t, 1, snapshot.Settings.BestOf)
}

func TestBestOfSeriesNeedsBothPlayersForNextRound(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, BestOf: 3, PlayerID: "one"})
	require.NoError(t, err)
	two, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, BestOf: 3, PlayerID: "two"})
	require.NoError(t, err)

	for _, move := range []struct {
		id string
		at int
	}{{one.ID, 0}, {two.ID, 3}, {one.ID, 1}, {two.ID, 4}, {one.ID, 2}} {
		require.NoError(t, room.HandleMove(move.id, move.at))
	}
	require.Equal(t, "round-finished", room.Snapshot().Status)
	require.NoError(t, room.ReadyRematch(one.ID))
	require.Equal(t, "round-finished", room.Snapshot().Status)
	require.NoError(t, room.ReadyRematch(two.ID))
	require.Equal(t, "playing", room.Snapshot().Status)
	require.Equal(t, 2, room.Snapshot().Round)
	require.Equal(t, "two", room.Snapshot().Turn)
}

func TestSpectatorsCanChatButCannotMove(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.Join(TicTacToe, "one", "One")
	require.NoError(t, err)
	_, room, err := hub.Join(TicTacToe, "two", "Two")
	require.NoError(t, err)
	spectator, spectatorRoom, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, Role: "spectator", RoomID: room.id, PlayerID: "watcher", PlayerName: "Watcher"})
	require.NoError(t, err)
	require.Same(t, room, spectatorRoom)
	require.EqualError(t, room.HandleMove(spectator.ID, 0), "player is not seated in this room")
	require.NoError(t, room.AddChat(spectator.ID, "  good luck!  "))
	snapshot := room.Snapshot()
	require.Equal(t, 1, snapshot.Spectators)
	require.Equal(t, "good luck!", snapshot.Chat[0].Content)
	room.mu.Lock()
	require.Equal(t, "one", room.snapshotForLocked(one).ViewerMark)
	require.Equal(t, "spectator", room.snapshotForLocked(spectator).ViewerMark)
	room.mu.Unlock()
	require.NoError(t, room.HandleMove(one.ID, 0))
}

func TestChatRejectsFloodsAndControlCharacters(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.Join(TicTacToe, "one", "One")
	require.NoError(t, err)
	_, room, err := hub.Join(TicTacToe, "two", "Two")
	require.NoError(t, err)
	require.NoError(t, room.AddChat(one.ID, "hi\u0000there"))
	require.Equal(t, "hithere", room.Snapshot().Chat[0].Content)
	require.EqualError(t, room.AddChat(one.ID, "again"), "messages are being sent too quickly")
}

func TestStartingAnotherGameCannotReuseAStaleRoom(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.Join(TicTacToe, "same-player", "One")
	require.NoError(t, err)
	_, oldRoom, err := hub.Join(TicTacToe, "old-rival", "Old Rival")
	require.NoError(t, err)
	require.Same(t, oldRoom, hub.playerRooms[one.ID])

	newPlayer, room, err := hub.Join(ConnectFour, "same-player", "One")
	require.NoError(t, err)
	require.Nil(t, room)
	require.Nil(t, hub.playerRooms[newPlayer.ID])
	_, newRoom, err := hub.Join(ConnectFour, "new-rival", "New Rival")
	require.NoError(t, err)
	require.NotNil(t, newRoom)
	require.Equal(t, ConnectFour, newRoom.Snapshot().Kind)
}

func TestIntentionalLeaveReleasesPlayerIdentity(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.Join(TicTacToe, "one", "One")
	require.NoError(t, err)
	_, room, err := hub.Join(TicTacToe, "two", "Two")
	require.NoError(t, err)

	hub.LeaveRoom(one, room)
	require.Nil(t, hub.playerRooms[one.ID])
	require.Equal(t, "finished", room.Snapshot().Status)
	require.Equal(t, "two", room.Snapshot().Winner)
}

func TestChessRoomEnforcesColorAndLegalMoves(t *testing.T) {
	hub := NewHub()
	white, _, err := hub.Join(Chess, "white", "White")
	require.NoError(t, err)
	black, room, err := hub.Join(Chess, "black", "Black")
	require.NoError(t, err)
	require.EqualError(t, room.HandleChessMove(black.ID, "e7", "e5", ""), "it is not your turn")
	require.NoError(t, room.HandleChessMove(white.ID, "e2", "e4", ""))
	require.Equal(t, "two", room.Snapshot().Turn)
	require.Equal(t, "e4", room.Snapshot().Chess.Moves[0].SAN)
}

func TestChessClockBanksApplyIncrement(t *testing.T) {
	hub := NewHub()
	white, _, err := hub.JoinMatch(JoinOptions{Kind: Chess, PlayerID: "white", Mode: "blitz"})
	require.NoError(t, err)
	_, room, err := hub.JoinMatch(JoinOptions{Kind: Chess, PlayerID: "black", Mode: "blitz"})
	require.NoError(t, err)
	before := room.Snapshot().Clocks["one"]
	require.LessOrEqual(t, before, int64(180000))
	// Opening moves use a separate window; the bank and increment start after.
	require.NoError(t, room.HandleChessMove(white.ID, "e2", "e4", ""))
	require.NoError(t, room.HandleChessMove("black", "e7", "e5", ""))
	require.NoError(t, room.HandleChessMove(white.ID, "g1", "f3", ""))
	after := room.Snapshot().Clocks
	require.Greater(t, after["one"], int64(181000))
	require.LessOrEqual(t, after["two"], int64(180000))
}

func TestNewTableRequestReleasesFinishedRoomButNotLiveGame(t *testing.T) {
	hub := NewHub()
	one, _, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "one", InviteCode: "TAAAAA"})
	require.NoError(t, err)
	_, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "two", InviteCode: "TAAAAA"})
	require.NoError(t, err)

	// Mid-game, asking for another table resumes the live game instead.
	_, resumed, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "one", InviteCode: "TBBBBB"})
	require.NoError(t, err)
	require.Same(t, room, resumed)

	require.NoError(t, room.Resign(one.ID))
	// Same table after the game: reconnect.
	_, again, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "one", InviteCode: "TAAAAA"})
	require.NoError(t, err)
	require.Same(t, room, again)
	// Different table after the game: queue for the new one.
	_, fresh, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "one", InviteCode: "TBBBBB"})
	require.NoError(t, err)
	require.Nil(t, fresh)
	require.False(t, room.Snapshot().Players[0].Connected)
	require.True(t, room.Snapshot().Players[0].Left)
}

func TestAbandonedSeriesIsRemoved(t *testing.T) {
	previous := roomRetention
	roomRetention = 50 * time.Millisecond
	defer func() { roomRetention = previous }()

	hub := NewHub()
	one, _, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "a1", BestOf: 3})
	require.NoError(t, err)
	two, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "a2", BestOf: 3})
	require.NoError(t, err)
	require.NoError(t, room.Resign(one.ID)) // round over, series still open

	hub.Disconnect(one, room, one.generation.Load())
	hub.Disconnect(two, room, two.generation.Load())
	require.Eventually(t, func() bool { return hub.Stats().Rooms == 0 && hub.Stats().Connections == 0 }, time.Second, 10*time.Millisecond)
}

func TestRoomWithAConnectedPlayerIsKept(t *testing.T) {
	previous := roomRetention
	roomRetention = 50 * time.Millisecond
	defer func() { roomRetention = previous }()

	hub := NewHub()
	one, _, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "k1", BestOf: 3})
	require.NoError(t, err)
	_, room, err := hub.JoinMatch(JoinOptions{Kind: TicTacToe, PlayerID: "k2", BestOf: 3})
	require.NoError(t, err)
	hub.Disconnect(one, room, one.generation.Load())
	time.Sleep(150 * time.Millisecond)
	require.Equal(t, 1, hub.Stats().Rooms)
}
