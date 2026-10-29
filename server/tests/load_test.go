package tests

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"boredgamz/core/arena"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestConcurrentWebSocketMatches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concurrent socket test in short mode")
	}
	testServer := newTestServer(t)
	defer testServer.Close()
	base := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/api/v1/arena/tictactoe"
	errors := make(chan error, 24)
	var wait sync.WaitGroup
	for match := 0; match < 24; match++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			errors <- runSocketMatch(base, index)
		}(match)
	}
	wait.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
}

func runSocketMatch(base string, index int) error {
	header := http.Header{"Origin": []string{"http://allowed.example"}}
	code := fmt.Sprintf("LOAD%02d", index)
	firstURL := fmt.Sprintf("%s?playerId=load-%d-a&name=Alice&inviteCode=%s", base, index, code)
	secondURL := fmt.Sprintf("%s?playerId=load-%d-b&name=Bob&inviteCode=%s", base, index, code)
	first, _, err := websocket.DefaultDialer.Dial(firstURL, header)
	if err != nil {
		return err
	}
	defer first.Close()
	_ = first.SetReadDeadline(time.Now().Add(8 * time.Second))
	var waiting arena.ServerMessage
	if err := first.ReadJSON(&waiting); err != nil {
		return err
	}
	if waiting.Type != "waiting" {
		return fmt.Errorf("match %d did not enter waiting", index)
	}
	second, _, err := websocket.DefaultDialer.Dial(secondURL, header)
	if err != nil {
		return err
	}
	defer second.Close()
	_ = second.SetReadDeadline(time.Now().Add(8 * time.Second))
	if _, err := loadReadState(first); err != nil {
		return err
	}
	if _, err := loadReadState(second); err != nil {
		return err
	}
	for _, move := range []struct {
		connection *websocket.Conn
		position   int
	}{{first, 0}, {second, 3}, {first, 1}, {second, 4}, {first, 2}} {
		if err := move.connection.WriteJSON(arena.ClientMessage{Type: "move", Position: move.position}); err != nil {
			return err
		}
		left, err := loadReadState(first)
		if err != nil {
			return err
		}
		right, err := loadReadState(second)
		if err != nil {
			return err
		}
		if fmt.Sprint(left.Board) != fmt.Sprint(right.Board) {
			return fmt.Errorf("match %d clients diverged", index)
		}
		if move.position == 2 && (left.Status != "finished" || left.Winner != "one") {
			return fmt.Errorf("match %d did not finish correctly", index)
		}
	}
	return nil
}

func loadReadState(connection *websocket.Conn) (*arena.Snapshot, error) {
	for attempt := 0; attempt < 5; attempt++ {
		var message arena.ServerMessage
		if err := connection.ReadJSON(&message); err != nil {
			return nil, err
		}
		if message.Type == "state" && message.State != nil {
			return message.State, nil
		}
	}
	return nil, fmt.Errorf("state message was not received")
}
