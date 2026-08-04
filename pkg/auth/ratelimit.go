package auth

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Limiter is a fixed window counter keyed by an arbitrary string.
// It is in-process on purpose: the application is deployed as a single instance,
// and a shared store would be needed only once that stops being true.
type Limiter struct {
	mu      sync.Mutex
	hits    map[string]*window
	limit   int
	period  time.Duration
	lastGC  time.Time
	gcEvery time.Duration
}

type window struct {
	count     int
	expiresAt time.Time
}

// NewLimiter allows limit attempts per key within period.
func NewLimiter(limit int, period time.Duration) *Limiter {
	return &Limiter{
		hits:    make(map[string]*window),
		limit:   limit,
		period:  period,
		gcEvery: period,
	}
}

// Allow records an attempt and reports whether it stays within the limit.
func (l *Limiter) Allow(key string) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.collect(now)

	w, ok := l.hits[key]
	if !ok || now.After(w.expiresAt) {
		l.hits[key] = &window{count: 1, expiresAt: now.Add(l.period)}
		return true
	}

	w.count++
	return w.count <= l.limit
}

// Reset drops the counter for a key, called after a successful attempt so a user
// who eventually types the right password is not punished for the earlier tries.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.hits, key)
}

// collect removes expired windows so the map cannot grow without bound.
// The caller must hold the lock.
func (l *Limiter) collect(now time.Time) {
	if now.Sub(l.lastGC) < l.gcEvery {
		return
	}
	l.lastGC = now
	for key, w := range l.hits {
		if now.After(w.expiresAt) {
			delete(l.hits, key)
		}
	}
}

// ClientIP returns the address to rate limit against.
// The leftmost X-Forwarded-For entry is deliberately ignored: a browser can set
// that header itself, which would let an attacker pick a fresh bucket per request.
// Only a platform header injected by a trusted proxy, or the socket address, is used.
func ClientIP(r *http.Request) string {
	for _, header := range []string{"X-Vercel-Forwarded-For", "CF-Connecting-IP", "X-Real-IP"} {
		if value := strings.TrimSpace(r.Header.Get(header)); value != "" {
			return value
		}
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
