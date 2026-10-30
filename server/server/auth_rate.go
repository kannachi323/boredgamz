package server

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type authRateWindow struct {
	count int
	reset time.Time
}

type authRateLimiter struct {
	mu      sync.Mutex
	windows map[string]authRateWindow
	limit   int
	window  time.Duration
	message string
	now     func() time.Time
}

func newAuthRateLimiter(limit int, window time.Duration) *authRateLimiter {
	return newRateLimiter(limit, window, "too many authentication attempts")
}

// newRateLimiter is a fixed-window, per-client-IP limiter. State is in memory,
// which matches the single-instance deployment model.
func newRateLimiter(limit int, window time.Duration, message string) *authRateLimiter {
	return &authRateLimiter{windows: map[string]authRateWindow{}, limit: limit, window: window, message: message, now: time.Now}
}

func (limiter *authRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, retryAfter := limiter.allow(clientAddress(r.RemoteAddr))
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(max(1, int(retryAfter.Seconds()))))
			http.Error(w, limiter.message, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (limiter *authRateLimiter) allow(key string) (bool, time.Duration) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	now := limiter.now()
	current := limiter.windows[key]
	if current.reset.IsZero() || !now.Before(current.reset) {
		limiter.windows[key] = authRateWindow{count: 1, reset: now.Add(limiter.window)}
		limiter.pruneLocked(now)
		return true, 0
	}
	if current.count >= limiter.limit {
		return false, current.reset.Sub(now)
	}
	current.count++
	limiter.windows[key] = current
	return true, 0
}

func (limiter *authRateLimiter) pruneLocked(now time.Time) {
	if len(limiter.windows) < 10_000 {
		return
	}
	for key, window := range limiter.windows {
		if !now.Before(window.reset) {
			delete(limiter.windows, key)
		}
	}
}

func clientAddress(remote string) string {
	host, _, err := net.SplitHostPort(remote)
	if err == nil && host != "" {
		return host
	}
	if remote == "" {
		return "unknown"
	}
	return remote
}
