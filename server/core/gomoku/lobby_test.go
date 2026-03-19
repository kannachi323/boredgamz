package gomoku

import (
	"boredgamz/core"
	"boredgamz/db"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func newTestServerConn(t *testing.T) (*websocket.Conn, func()) {
	t.Helper()

	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	serverConnCh := make(chan *websocket.Conn, 1)
	errCh := make(chan error, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			errCh <- err
			return
		}
		serverConnCh <- conn
	}))

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second):
	}

	var serverConn *websocket.Conn
	select {
	case serverConn = <-serverConnCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for test websocket server connection")
	}

	cleanup := func() {
		_ = clientConn.Close()
		_ = serverConn.Close()
		ts.Close()
	}

	return serverConn, cleanup
}

func newLobbyTestPlayer(playerID string, color string, conn *websocket.Conn) *core.Player {
	return core.NewPlayer(
		playerID,
		"player-"+playerID,
		color,
		core.NewPlayerClock("Rapid"),
		conn,
	)
}

func TestGomokuLobby_DeduplicatesPlayerIDInQueue(t *testing.T) {
	controller := NewGomokuLobby(100, "gomoku-casual-9x9", &db.Database{})
	lobby, ok := controller.(*GomokuLobby)
	require.True(t, ok)

	conn1, cleanup1 := newTestServerConn(t)
	defer cleanup1()

	player1 := newLobbyTestPlayer("same-player", "black", conn1)
	lobby.AddPlayer(player1)

	require.Equal(t, 1, lobby.BlackQueue.Len())
	require.Equal(t, 1, lobby.NumPlayers)

	conn2, cleanup2 := newTestServerConn(t)
	defer cleanup2()

	duplicate := newLobbyTestPlayer("same-player", "black", conn2)
	lobby.AddPlayer(duplicate)

	require.Equal(t, 1, lobby.BlackQueue.Len(), "duplicate join should not add extra queue entry")
	require.Equal(t, 1, lobby.NumPlayers, "duplicate join should not increase player count")

	slot, exists := lobby.PlayerSlot["same-player"]
	require.True(t, exists)
	require.NotNil(t, slot)
	require.NotNil(t, slot.Player)
	require.Same(t, conn2, slot.Player.Conn, "duplicate join should refresh to the newest connection")
}

func TestGomokuLobby_SelfHealsDisconnectedQueueEntries(t *testing.T) {
	controller := NewGomokuLobby(100, "gomoku-casual-9x9", &db.Database{})
	lobby, ok := controller.(*GomokuLobby)
	require.True(t, ok)

	staleConn, staleCleanup := newTestServerConn(t)
	defer staleCleanup()

	stale := newLobbyTestPlayer("stale-player", "black", staleConn)
	lobby.AddPlayer(stale)
	require.Equal(t, 1, lobby.BlackQueue.Len())

	stale.ClosePlayer()
	require.Eventually(t, func() bool {
		return stale.Disconnected.Load()
	}, 2*time.Second, 50*time.Millisecond)

	freshConn, freshCleanup := newTestServerConn(t)
	defer freshCleanup()

	fresh := newLobbyTestPlayer("fresh-player", "black", freshConn)
	lobby.AddPlayer(fresh)

	require.Equal(t, 1, lobby.BlackQueue.Len(), "stale queue entry should be removed automatically")
	require.Equal(t, 1, lobby.NumPlayers, "lobby should self-restore count after stale removal")
	_, staleExists := lobby.PlayerSlot["stale-player"]
	require.False(t, staleExists)
	_, freshExists := lobby.PlayerSlot["fresh-player"]
	require.True(t, freshExists)
}

func TestGomokuLobby_MatchesPlayersAndRegistersRoom(t *testing.T) {
	controller := NewGomokuLobby(100, "gomoku-casual-9x9", &db.Database{})
	lobby, ok := controller.(*GomokuLobby)
	require.True(t, ok)

	blackConn, blackCleanup := newTestServerConn(t)
	defer blackCleanup()
	whiteConn, whiteCleanup := newTestServerConn(t)
	defer whiteCleanup()

	black := newLobbyTestPlayer("black-player", "black", blackConn)
	white := newLobbyTestPlayer("white-player", "white", whiteConn)

	lobby.AddPlayer(black)
	lobby.AddPlayer(white)

	require.Eventually(t, func() bool {
		_, ok1 := lobby.RoomManager.GetPlayerRoom("black-player")
		_, ok2 := lobby.RoomManager.GetPlayerRoom("white-player")
		return ok1 && ok2
	}, 3*time.Second, 50*time.Millisecond)

	require.Equal(t, 0, lobby.BlackQueue.Len())
	require.Equal(t, 0, lobby.WhiteQueue.Len())
}
