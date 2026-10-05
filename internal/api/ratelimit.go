package api

import (
	"sync"
	"time"
)

// rateLimiter allows each client a fixed number of requests per time window
// (a fixed-window counter). It's in memory, so limits reset when the server
// restarts and aren't shared between multiple server instances.
type rateLimiter struct {
	mu        sync.Mutex
	limit     int
	window    time.Duration
	clients   map[string]*clientWindow
	lastSweep time.Time
	now       func() time.Time // swapped out in tests
}

type clientWindow struct {
	start time.Time
	count int
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientWindow),
		now:     time.Now,
	}
}

// allow records a request from key and reports whether it's within the limit.
// When it isn't, retryAfter says how long until that client's window resets.
func (l *rateLimiter) allow(key string) (ok bool, retryAfter time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.sweep(now)

	c, found := l.clients[key]
	if !found || now.Sub(c.start) >= l.window {
		l.clients[key] = &clientWindow{start: now, count: 1}
		return true, 0
	}
	if c.count >= l.limit {
		return false, c.start.Add(l.window).Sub(now)
	}
	c.count++
	return true, 0
}

// sweep drops expired windows so the map doesn't grow forever. It runs at most
// once per window, which spreads the cost across requests.
func (l *rateLimiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < l.window {
		return
	}
	for key, c := range l.clients {
		if now.Sub(c.start) >= l.window {
			delete(l.clients, key)
		}
	}
	l.lastSweep = now
}
