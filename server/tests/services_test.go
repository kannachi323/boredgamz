package tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// End-to-end checks over real WebSockets for every socket service.

type wireMessage struct {
	Type    string          `json:"type"`
	Message string          `json:"message"`
	State   json.RawMessage `json:"state"`
}

func dial(t *testing.T, base, path string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(base, "http")+path, http.Header{"Origin": []string{"http://allowed.example"}})
	require.NoError(t, err)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	t.Cleanup(func() { conn.Close() })
	return conn
}

func next(t *testing.T, conn *websocket.Conn, kind string) wireMessage {
	t.Helper()
	for attempt := 0; attempt < 6; attempt++ {
		var message wireMessage
		require.NoError(t, conn.ReadJSON(&message))
		if message.Type == kind {
			return message
		}
	}
	t.Fatalf("no %q message", kind)
	return wireMessage{}
}

func TestPokerOverSockets(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()
	a := dial(t, server.URL, "/api/v1/poker?playerId=pa&name=Ada&mode=holdem")
	require.Equal(t, "waiting", next(t, a, "waiting").Type)
	b := dial(t, server.URL, "/api/v1/poker?playerId=pb&name=Bo&mode=holdem")

	type seat struct {
		Cards     []json.RawMessage `json:"cards"`
		CardCount int               `json:"cardCount"`
	}
	type pokerState struct {
		ViewerSeat int     `json:"viewerSeat"`
		Turn       int     `json:"turn"`
		Phase      string  `json:"phase"`
		Seats      [2]seat `json:"seats"`
	}
	var stateA, stateB pokerState
	require.NoError(t, json.Unmarshal(next(t, a, "state").State, &stateA))
	require.NoError(t, json.Unmarshal(next(t, b, "state").State, &stateB))
	require.Equal(t, 0, stateA.ViewerSeat)
	require.Equal(t, 1, stateB.ViewerSeat)
	// Each player sees only their own hole cards.
	require.Len(t, stateA.Seats[0].Cards, 2)
	require.Empty(t, stateA.Seats[1].Cards)
	require.Equal(t, 2, stateA.Seats[1].CardCount)
	require.Len(t, stateB.Seats[1].Cards, 2)
	require.Empty(t, stateB.Seats[0].Cards)

	// The wrong seat can't act; the right one folds and the hand ends.
	actor, other := a, b
	if stateA.Turn == 1 {
		actor, other = b, a
	}
	require.NoError(t, other.WriteJSON(map[string]any{"type": "action", "action": "fold"}))
	require.Equal(t, "error", next(t, other, "error").Type)
	require.NoError(t, actor.WriteJSON(map[string]any{"type": "action", "action": "fold"}))
	var after pokerState
	require.NoError(t, json.Unmarshal(next(t, other, "state").State, &after))
	require.Equal(t, "hand-finished", after.Phase)
}

func TestMahjongOverSockets(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()
	a := dial(t, server.URL, "/api/v1/mahjong?playerId=ma&name=Ada&mode=duel")
	require.Equal(t, "waiting", next(t, a, "waiting").Type)
	b := dial(t, server.URL, "/api/v1/mahjong?playerId=mb&name=Bo&mode=duel")

	type mahjongState struct {
		ViewerSeat int    `json:"viewerSeat"`
		Status     string `json:"status"`
		Free       []int  `json:"free"`
		Tiles      []struct {
			Face string `json:"face"`
		} `json:"tiles"`
	}
	var stateA, stateB mahjongState
	require.NoError(t, json.Unmarshal(next(t, a, "state").State, &stateA))
	require.NoError(t, json.Unmarshal(next(t, b, "state").State, &stateB))
	require.Equal(t, 0, stateA.ViewerSeat)
	require.Equal(t, 1, stateB.ViewerSeat)
	require.Equal(t, "playing", stateA.Status)
	require.Len(t, stateA.Tiles, 72)

	// A tile can't pair with itself.
	require.NoError(t, a.WriteJSON(map[string]any{"type": "pair", "first": stateA.Free[0], "second": stateA.Free[0]}))
	require.Equal(t, "error", next(t, a, "error").Type)

	// A hint names a real free pair; claiming it is accepted.
	require.NoError(t, a.WriteJSON(map[string]any{"type": "hint"}))
	var hinted struct {
		Hint []int `json:"hint"`
	}
	require.NoError(t, json.Unmarshal(next(t, a, "state").State, &hinted))
	require.Len(t, hinted.Hint, 2)
	require.NoError(t, a.WriteJSON(map[string]any{"type": "pair", "first": hinted.Hint[0], "second": hinted.Hint[1]}))
	var claimed struct {
		Removed []bool `json:"removed"`
	}
	require.NoError(t, json.Unmarshal(next(t, a, "state").State, &claimed))
	require.True(t, claimed.Removed[hinted.Hint[0]])
	require.True(t, claimed.Removed[hinted.Hint[1]])
}

func TestChessOverSockets(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()
	white := dial(t, server.URL, "/api/v1/arena/chess?playerId=cw&name=Ada&mode=blitz")
	require.Equal(t, "waiting", next(t, white, "waiting").Type)
	black := dial(t, server.URL, "/api/v1/arena/chess?playerId=cb&name=Bo&mode=blitz")
	next(t, white, "state")
	next(t, black, "state")

	require.NoError(t, white.WriteJSON(map[string]any{"type": "move", "from": "e2", "to": "e5"}))
	require.Equal(t, "error", next(t, white, "error").Type)
	require.NoError(t, white.WriteJSON(map[string]any{"type": "move", "from": "e2", "to": "e4"}))
	var state struct {
		Turn  string `json:"turn"`
		Chess struct {
			Moves []struct {
				SAN string `json:"san"`
			} `json:"moves"`
		} `json:"chess"`
		Settings struct {
			TimeControl string `json:"timeControl"`
		} `json:"settings"`
	}
	require.NoError(t, json.Unmarshal(next(t, black, "state").State, &state))
	require.Equal(t, "two", state.Turn)
	require.Equal(t, "e4", state.Chess.Moves[0].SAN)
	require.Equal(t, "3m+2s", state.Settings.TimeControl)
}

func TestCancelledQueueEntryIsNeverMatched(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()
	ghost := dial(t, server.URL, "/api/v1/arena/gomoku?playerId=ghost&name=Ghost")
	require.Equal(t, "waiting", next(t, ghost, "waiting").Type)
	ghost.Close()

	waitingCount := func() int {
		response, err := http.Get(server.URL + "/api/v1/status")
		require.NoError(t, err)
		defer response.Body.Close()
		var status struct {
			Arena struct {
				Waiting int `json:"waiting"`
			} `json:"arena"`
		}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&status))
		return status.Arena.Waiting
	}
	require.Eventually(t, func() bool { return waitingCount() == 0 }, 2*time.Second, 20*time.Millisecond)

	// The next player waits instead of being seated against the ghost.
	real := dial(t, server.URL, "/api/v1/arena/gomoku?playerId=real&name=Ada")
	require.Equal(t, "waiting", next(t, real, "waiting").Type)
}
