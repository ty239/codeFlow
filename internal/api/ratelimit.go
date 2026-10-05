package api

import (
	"math"
	"net"
	"net/http"
	"strconv"
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

// wrap returns a handler that rejects requests over the limit with
// 429 Too Many Requests and a Retry-After header, keyed by client IP.
func (l *rateLimiter) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, retryAfter := l.allow(clientIP(r))
		if !ok {
			w.Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(retryAfter.Seconds()))))
			writeError(w, http.StatusTooManyRequests, "too many requests, please try again later")
			return
		}
		next(w, r)
	}
}

// clientIP returns the address of the directly connected client. It ignores
// X-Forwarded-For on purpose: any client can set that header to dodge the
// limit. If the server is later put behind a trusted reverse proxy, read the
// proxy's header here instead.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
