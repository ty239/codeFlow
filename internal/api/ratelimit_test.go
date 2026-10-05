package api

import (
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
