package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	appserver "boredgamz/server"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// TestAccountsAndAuthenticatedSockets runs against a real Postgres when
// TEST_DATABASE_URL is set (CI and local); it is skipped otherwise.
func TestAccountsAndAuthenticatedSockets(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	t.Setenv("APP_ENV", "test")
	t.Setenv("DATABASE_URL", dsn)
	t.Setenv("JWT_SECRET_KEY", "integration-test-secret-0123456789abcdef")
	t.Setenv("ALLOWED_ORIGINS", "http://allowed.example")
	server, err := appserver.CreateServer()
	require.NoError(t, err)
	defer server.DB.Stop()
	httpServer := httptest.NewServer(server.Router)
	defer httpServer.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	email := fmt.Sprintf("it-%d@example.com", time.Now().UnixNano())
	post := func(path string, body any) *http.Response {
		payload, _ := json.Marshal(body)
		request, _ := http.NewRequest("POST", httpServer.URL+path, bytes.NewReader(payload))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", "http://allowed.example")
		response, err := client.Do(request)
		require.NoError(t, err)
		response.Body.Close()
		return response
	}

	require.Equal(t, http.StatusBadRequest, post("/api/v1/auth/signup", map[string]string{"username": "Ada", "email": email, "password": strings.Repeat("x", 73)}).StatusCode)
	require.Equal(t, http.StatusCreated, post("/api/v1/auth/signup", map[string]string{"username": "Ada", "email": email, "password": "correct horse battery"}).StatusCode)
	require.Equal(t, http.StatusConflict, post("/api/v1/auth/signup", map[string]string{"username": "Ada", "email": email, "password": "correct horse battery"}).StatusCode)
	require.Equal(t, http.StatusUnauthorized, post("/api/v1/auth/login", map[string]string{"email": email, "password": "wrong password!!"}).StatusCode)
	require.Equal(t, http.StatusOK, post("/api/v1/auth/login", map[string]string{"email": email, "password": "correct horse battery"}).StatusCode)

	session, err := client.Get(httpServer.URL + "/api/v1/auth/session")
	require.NoError(t, err)
	var me struct {
		Username string `json:"username"`
		ID       string `json:"id"`
	}
	require.NoError(t, json.NewDecoder(session.Body).Decode(&me))
	session.Body.Close()
	require.Equal(t, "Ada", me.Username)

	wsURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/api/v1/arena/tictactoe?playerId=browser:tab1&name=Impostor"
	header := http.Header{"Origin": []string{"http://allowed.example"}}

	// No cookie: rejected before the upgrade.
	_, response, err := websocket.DefaultDialer.Dial(wsURL, header)
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	// With the session cookie the account name wins over the query string.
	base, _ := url.Parse(httpServer.URL)
	for _, cookie := range jar.Cookies(base) {
		header.Add("Cookie", cookie.Name+"="+cookie.Value)
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	require.NoError(t, err)
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	require.Equal(t, "waiting", readMessage(t, conn).Type)

	// Cross-site logout is refused; same-origin logout clears the session.
	request, _ := http.NewRequest("POST", httpServer.URL+"/api/v1/auth/logout", nil)
	request.Header.Set("Origin", "http://evil.example")
	response, err = client.Do(request)
	require.NoError(t, err)
	response.Body.Close()
	require.Equal(t, http.StatusForbidden, response.StatusCode)
	require.Equal(t, http.StatusNoContent, post("/api/v1/auth/logout", nil).StatusCode)
	session, err = client.Get(httpServer.URL + "/api/v1/auth/session")
	require.NoError(t, err)
	session.Body.Close()
	require.Equal(t, http.StatusUnauthorized, session.StatusCode)
}
