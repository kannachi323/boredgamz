package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"playgomoku/backend/api"
	"playgomoku/backend/game"
	"playgomoku/backend/server"
)

func executeRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

//new game state
func TestNewGameState(t *testing.T) {
	s := server.CreateServer()
	s.MountHandlers()

	p1 := game.Player{
		ID: "1",
		Username: "Alice",
	}

	p2 := game.Player{
		ID: "2",
		Username: "Bob",
	}

	reqBody := api.NewGameStateRequest{
		Size: 19,
		P1: p1,
		P2: p2,
	}

	jsonBytes, _ := json.Marshal(reqBody)


	req, err := http.NewRequest("POST", "/new-game-state", bytes.NewReader(jsonBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, response.Code)

	var actual *game.GameState
	err = json.Unmarshal(response.Body.Bytes(), &actual)
	require.NoError(t, err)

	expected := game.CreateGameState(19, &p1, &p2)
	

	require.Equal(t, expected, actual)
}