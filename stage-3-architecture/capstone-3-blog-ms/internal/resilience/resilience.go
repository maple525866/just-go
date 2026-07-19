package resilience

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sony/gobreaker/v2"
)

var ErrLimited = errors.New("rate limit exceeded")

type Limiter struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64
	last     time.Time
}

func NewLimiter(capacity int, refillPerSecond float64) *Limiter {
	now := time.Now()
	return &Limiter{capacity: float64(capacity), tokens: float64(capacity), rate: refillPerSecond, last: now}
}

func (l *Limiter) Allow() bool {
	if l == nil {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.tokens += now.Sub(l.last).Seconds() * l.rate
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.last = now
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

type Breaker[T any] struct {
	breaker *gobreaker.CircuitBreaker[T]
}

func NewBreaker[T any](name string, failures uint32, timeout time.Duration) *Breaker[T] {
	if failures == 0 {
		failures = 3
	}
	return &Breaker[T]{breaker: gobreaker.NewCircuitBreaker[T](gobreaker.Settings{
		Name: name, MaxRequests: 1, Timeout: timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= failures },
	})}
}

func (b *Breaker[T]) Execute(operation func() (T, error)) (T, error) {
	if b == nil {
		return operation()
	}
	return b.breaker.Execute(operation)
}

func Retry[T any](ctx context.Context, attempts int, delay time.Duration, retryable func(error) bool, operation func(context.Context) (T, error)) (T, error) {
	var zero T
	if attempts < 1 {
		attempts = 1
	}
	for attempt := 1; attempt <= attempts; attempt++ {
		value, err := operation(ctx)
		if err == nil {
			return value, nil
		}
		if attempt == attempts || !retryable(err) {
			return zero, err
		}
		timer := time.NewTimer(delay * time.Duration(1<<(attempt-1)))
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, ctx.Err()
		case <-timer.C:
		}
	}
	return zero, nil
}
