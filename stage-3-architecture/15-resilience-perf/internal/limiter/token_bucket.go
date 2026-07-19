package limiter

import (
	"errors"
	"math"
	"sync"
	"time"
)

type Config struct {
	Capacity        int
	RefillPerSecond float64
}

type Decision struct {
	Allowed    bool
	RetryAfter time.Duration
	Remaining  int
}

type TokenBucket struct {
	mu              sync.Mutex
	capacity        float64
	refillPerSecond float64
	tokens          float64
	lastRefill      time.Time
	clock           Clock
}

func NewTokenBucket(config Config, clock Clock) (*TokenBucket, error) {
	if config.Capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	if config.RefillPerSecond <= 0 || math.IsNaN(config.RefillPerSecond) || math.IsInf(config.RefillPerSecond, 0) {
		return nil, errors.New("refill rate must be positive and finite")
	}
	if clock == nil {
		clock = realClock{}
	}
	now := clock.Now()
	return &TokenBucket{capacity: float64(config.Capacity), refillPerSecond: config.RefillPerSecond, tokens: float64(config.Capacity), lastRefill: now, clock: clock}, nil
}

func (b *TokenBucket) Allow() Decision {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refillLocked()
	if b.tokens >= 1 {
		b.tokens--
		return Decision{Allowed: true, Remaining: int(math.Floor(b.tokens))}
	}
	missing := 1 - b.tokens
	retryAfter := time.Duration(math.Ceil((missing / b.refillPerSecond) * float64(time.Second)))
	if retryAfter < 0 {
		retryAfter = 0
	}
	return Decision{Allowed: false, RetryAfter: retryAfter, Remaining: 0}
}

func (b *TokenBucket) refillLocked() {
	now := b.clock.Now()
	elapsed := now.Sub(b.lastRefill)
	if elapsed <= 0 {
		return
	}
	b.tokens += elapsed.Seconds() * b.refillPerSecond
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.lastRefill = now
}
