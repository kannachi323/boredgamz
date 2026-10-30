package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuthRateLimiterBlocksAndResetsPerClient(t *testing.T) {
	limiter := newAuthRateLimiter(2, time.Minute)
	now := time.Now()
	limiter.now = func() time.Time { return now }
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	for attempt, want := range []int{http.StatusNoContent, http.StatusNoContent, http.StatusTooManyRequests} {
		request := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		request.RemoteAddr = "203.0.113.10:4000"
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		require.Equal(t, want, recorder.Code, "attempt %d", attempt+1)
	}
	now = now.Add(61 * time.Second)
	request := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	request.RemoteAddr = "203.0.113.10:4001"
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestClientAddressNormalizesPorts(t *testing.T) {
	require.Equal(t, "203.0.113.10", clientAddress("203.0.113.10:4000"))
	require.Equal(t, "2001:db8::1", clientAddress("[2001:db8::1]:4000"))
}
