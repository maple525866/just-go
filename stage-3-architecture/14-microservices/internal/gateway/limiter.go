package gateway

import (
	"strings"
	"sync"
	"time"
)

type bucket struct {
	started time.Time
	window  time.Duration
	count   int
}

// Limiter is intentionally process-local. Production gateways generally need
// a distributed limiter when several instances share one quota.
type Limiter struct {
	mu      sync.Mutex
	clock   func() time.Time
	buckets map[string]bucket
}

func NewLimiter(clock func() time.Time) *Limiter {
	if clock == nil {
		clock = time.Now
	}
	return &Limiter{clock: clock, buckets: make(map[string]bucket)}
}

func (l *Limiter) Allow(key string, limit int, window time.Duration) bool {
	key = strings.TrimSpace(key)
	if key == "" || limit <= 0 || window <= 0 {
		return false
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock()
	current, exists := l.buckets[key]
	if !exists || current.window != window || !now.Before(current.started.Add(current.window)) {
		current = bucket{started: now, window: window}
	}
	if current.count >= limit {
		l.buckets[key] = current
		return false
	}
	current.count++
	l.buckets[key] = current
	return true
}
