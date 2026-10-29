package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"boredgamz/core/arena"
	appserver "boredgamz/server"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	t.Setenv("APP_ENV", "test")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ALLOWED_ORIGINS", "http://allowed.example")
	server, err := appserver.CreateServer()
	require.NoError(t, err)
	return httptest.NewServer(server.Router)
}

func TestConfiguredDatabaseFailureStopsServerInitialization(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("DATABASE_URL", "postgres://boredgamz:unavailable@127.0.0.1:1/boredgamz?sslmode=disable&connect_timeout=1")
	_, err := appserver.CreateServer()
	require.ErrorContains(t, err, "start account database")
}

func TestHealthAndCatalog(t *testing.T) {
	testServer := newTestServer(t)
	defer testServer.Close()

	response, err := http.Get(testServer.URL + "/healthz")
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "nosniff", response.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "same-origin", response.Header.Get("Cross-Origin-Opener-Policy"))
	require.Equal(t, "same-origin", response.Header.Get("Cross-Origin-Resource-Policy"))

	response, err = http.Get(testServer.URL + "/readyz")
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	response, err = http.Get(testServer.URL + "/api/v1/catalog")
	require.NoError(t, err)
	defer response.Body.Close()
	var body struct {
		Games []map[string]any `json:"games"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Len(t, body.Games, 6)

	response, err = http.Get(testServer.URL + "/api/v1/status")
	require.NoError(t, err)
	defer response.Body.Close()
	var status struct {
		Service string `json:"service"`
		Arena   struct {
			Waiting int `json:"waiting"`
			Rooms   int `json:"rooms"`
		} `json:"arena"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&status))
	require.Equal(t, "boredgamz", status.Service)
	require.Zero(t, status.Arena.Waiting)
	require.Zero(t, status.Arena.Rooms)
}

func TestArenaRejectsUnknownGamesAndUntrustedOrigins(t *testing.T) {
	testServer := newTestServer(t)
	defer testServer.Close()

	response, err := http.Get(testServer.URL + "/api/v1/arena/checkers?playerId=p1")
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/api/v1/arena/tictactoe?playerId=p1&name=Alice"
	header := http.Header{"Origin": []string{"http://evil.example"}}
	connection, response, err := websocket.DefaultDialer.Dial(wsURL, header)
	if connection != nil {
		connection.Close()
	}
	require.Error(t, err)
	require.NotNil(t, response)
	require.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestArenaMatchValidatesTurnsMovesAndWinner(t *testing.T) {
	testServer := newTestServer(t)
	defer testServer.Close()
	base := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/api/v1/arena/tictactoe"
	header := http.Header{"Origin": []string{"http://allowed.example"}}

	one, _, err := websocket.DefaultDialer.Dial(base+"?playerId=one&name=Alice", header)
	require.NoError(t, err)
	defer one.Close()
	one.SetReadDeadline(time.Now().Add(3 * time.Second))
	require.Equal(t, "waiting", readMessage(t, one).Type)

	two, _, err := websocket.DefaultDialer.Dial(base+"?playerId=two&name=Bob", header)
	require.NoError(t, err)
	defer two.Close()
	two.SetReadDeadline(time.Now().Add(3 * time.Second))
	firstState := readUntilState(t, two)
	require.Equal(t, "one", firstState.State.Turn)
	require.Equal(t, "playing", firstState.State.Status)
	readUntilState(t, one)

	require.NoError(t, two.WriteJSON(arena.ClientMessage{Type: "move", Position: 0}))
	require.Equal(t, "error", readMessage(t, two).Type)

	var final *arena.Snapshot
	for _, move := range []struct {
		connection *websocket.Conn
		position   int
	}{{one, 0}, {two, 3}, {one, 1}, {two, 4}, {one, 2}} {
		require.NoError(t, move.connection.WriteJSON(arena.ClientMessage{Type: "move", Position: move.position}))
		stateOne := readUntilState(t, one)
		stateTwo := readUntilState(t, two)
		require.Equal(t, stateOne.State.Board, stateTwo.State.Board)
		final = stateOne.State
	}

	require.NotNil(t, final)
	require.Equal(t, "finished", final.Status)
	require.Equal(t, "one", final.Winner)
	require.Equal(t, "one", final.Board[0])
	require.Equal(t, "one", final.Board[1])
	require.Equal(t, "one", final.Board[2])

	require.NoError(t, two.WriteJSON(arena.ClientMessage{Type: "move", Position: 8}))
	require.Equal(t, "error", readMessage(t, two).Type)
}

func readMessage(t *testing.T, connection *websocket.Conn) arena.ServerMessage {
	t.Helper()
	var message arena.ServerMessage
	require.NoError(t, connection.ReadJSON(&message))
	return message
}

func readUntilState(t *testing.T, connection *websocket.Conn) arena.ServerMessage {
	t.Helper()
	for attempt := 0; attempt < 4; attempt++ {
		message := readMessage(t, connection)
		if message.Type == "state" {
			return message
		}
	}
	t.Fatal("state message was not received")
	return arena.ServerMessage{}
}
