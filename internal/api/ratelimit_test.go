package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time { return c.t }

func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newTestLimiter(limit int, window time.Duration) (*rateLimiter, *fakeClock) {
	clock := &fakeClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	l := newRateLimiter(limit, window)
	l.now = clock.now
	return l, clock
}

func TestRateLimiterAllowsUpToLimit(t *testing.T) {
	l, _ := newTestLimiter(3, time.Minute)

	for i := range 3 {
		if ok, _ := l.allow("1.2.3.4"); !ok {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	ok, retryAfter := l.allow("1.2.3.4")
	if ok {
		t.Fatal("request over the limit should be blocked")
	}
	if retryAfter != time.Minute {
		t.Fatalf("expected retryAfter of 1m, got %v", retryAfter)
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	l, clock := newTestLimiter(1, time.Minute)

	l.allow("1.2.3.4")
	clock.advance(40 * time.Second)
	if ok, retryAfter := l.allow("1.2.3.4"); ok || retryAfter != 20*time.Second {
		t.Fatalf("expected block with 20s left, got ok=%v retryAfter=%v", ok, retryAfter)
	}

	clock.advance(20 * time.Second)
	if ok, _ := l.allow("1.2.3.4"); !ok {
		t.Fatal("request should be allowed once the window has passed")
	}
}

func TestRateLimiterTracksClientsSeparately(t *testing.T) {
	l, _ := newTestLimiter(1, time.Minute)

	l.allow("1.2.3.4")
	if ok, _ := l.allow("5.6.7.8"); !ok {
		t.Fatal("a different client should have its own limit")
	}
}

func TestRateLimiterSweepsExpiredClients(t *testing.T) {
	l, clock := newTestLimiter(5, time.Minute)

	l.allow("1.2.3.4")
	l.allow("5.6.7.8")
	clock.advance(2 * time.Minute)
	l.allow("9.9.9.9")

	if len(l.clients) != 1 {
		t.Fatalf("expected expired clients to be swept, %d remain", len(l.clients))
	}
}

func TestWrapReturns429WithRetryAfter(t *testing.T) {
	l, _ := newTestLimiter(2, time.Minute)
	h := l.wrap(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	codes := make([]int, 3)
	var last *httptest.ResponseRecorder
	for i := range codes {
		last = httptest.NewRecorder()
		h(last, httptest.NewRequest(http.MethodPost, "/login", nil))
		codes[i] = last.Code
	}

	if codes[0] != http.StatusOK || codes[1] != http.StatusOK || codes[2] != http.StatusTooManyRequests {
		t.Fatalf("expected [200 200 429], got %v", codes)
	}
	if got := last.Header().Get("Retry-After"); got != "60" {
		t.Fatalf("expected Retry-After 60, got %q", got)
	}
}

// The real routes must be wrapped. An invalid JSON body fails before any
// database access, so the allowed requests get 400 and the next one gets 429.
func TestAuthRoutesAreRateLimited(t *testing.T) {
	tests := []struct {
		path  string
		limit int
	}{
		{"/login", 5},
		{"/signup", 10},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			mux := http.NewServeMux()
			NewServer(nil, nil).RegisterRoutes(mux)

			for i := range tt.limit + 1 {
				w := httptest.NewRecorder()
				mux.ServeHTTP(w, httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader("{bad json")))

				want := http.StatusBadRequest
				if i == tt.limit {
					want = http.StatusTooManyRequests
				}
				if w.Code != want {
					t.Fatalf("request %d: expected status %d, got %d", i+1, want, w.Code)
				}
			}
		})
	}
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		remoteAddr string
		want       string
	}{
		{"192.0.2.1:1234", "192.0.2.1"},
		{"[2001:db8::1]:443", "2001:db8::1"},
		{"no-port", "no-port"},
	}

	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = tt.remoteAddr
		r.Header.Set("X-Forwarded-For", "10.0.0.1")
		if got := clientIP(r); got != tt.want {
			t.Errorf("clientIP(%q) = %q, want %q", tt.remoteAddr, got, tt.want)
		}
	}
}
