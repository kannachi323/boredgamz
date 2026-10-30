package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"boredgamz/db"

	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T, env string) *Server {
	t.Helper()
	t.Setenv("APP_ENV", env)
	return &Server{DB: &db.Database{}, allowedOrigins: parseAllowedOrigins("https://games.example")}
}

func TestTrustedOriginBlocksCrossSiteWrites(t *testing.T) {
	s := testServer(t, "production")
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	handler := s.requireTrustedOrigin(ok)

	cases := []struct {
		name, method, origin, referer string
		want                          int
	}{
		{"same origin post", "POST", "https://games.example", "", http.StatusNoContent},
		{"foreign origin post", "POST", "https://evil.example", "", http.StatusForbidden},
		{"referer fallback", "POST", "", "https://games.example/login", http.StatusNoContent},
		{"foreign referer", "POST", "", "https://evil.example/x", http.StatusForbidden},
		{"no origin in production", "POST", "", "", http.StatusForbidden},
		{"reads are not checked", "GET", "https://evil.example", "", http.StatusNoContent},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			request := httptest.NewRequest(c.method, "/v1/auth/logout", nil)
			if c.origin != "" {
				request.Header.Set("Origin", c.origin)
			}
			if c.referer != "" {
				request.Header.Set("Referer", c.referer)
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			require.Equal(t, c.want, recorder.Code)
		})
	}
}

func TestSocketsRequireAccountWhenAccountsEnabled(t *testing.T) {
	s := testServer(t, "test")
	// Guest mode: no database configured.
	request := httptest.NewRequest("GET", "/v1/poker?playerId=abc&name=Guest", nil)
	params, ok := s.prepareSocket(httptest.NewRecorder(), request)
	require.True(t, ok)
	require.Equal(t, "abc", params.PlayerID)

	// Invalid role is rejected before identity is considered.
	request = httptest.NewRequest("GET", "/v1/poker?playerId=abc&role=admin", nil)
	recorder := httptest.NewRecorder()
	_, ok = s.prepareSocket(recorder, request)
	require.False(t, ok)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestRequestLoggerPassesThroughStatus(t *testing.T) {
	handler := requestLogger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) }))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/api/v1/catalog", nil))
	require.Equal(t, http.StatusTeapot, recorder.Code)
}
