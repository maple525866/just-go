package retry

import (
	"context"
	"time"
)

type Policy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      func(time.Duration, int) time.Duration
	Sleep       func(context.Context, time.Duration) error
	ShouldRetry func(error) bool
}

type Stats struct {
	Attempts int
	Delays   []time.Duration
}

func Do[T any](ctx context.Context, policy Policy, op func(context.Context) (T, error)) (T, Stats, error) {
	var zero T
	if ctx == nil {
		ctx = context.Background()
	}
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if policy.BaseDelay <= 0 {
		policy.BaseDelay = time.Millisecond
	}
	if policy.MaxDelay <= 0 {
		policy.MaxDelay = policy.BaseDelay
	}
	if policy.Sleep == nil {
		policy.Sleep = sleepContext
	}
	if policy.ShouldRetry == nil {
		policy.ShouldRetry = func(error) bool { return false }
	}

	stats := Stats{}
	delay := policy.BaseDelay
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return zero, stats, err
		}
		stats.Attempts++
		result, err := op(ctx)
		if err == nil {
			return result, stats, nil
		}
		if attempt == policy.MaxAttempts || !policy.ShouldRetry(err) {
			return zero, stats, err
		}
		sleepFor := delay
		if sleepFor > policy.MaxDelay {
			sleepFor = policy.MaxDelay
		}
		if policy.Jitter != nil {
			sleepFor = policy.Jitter(sleepFor, attempt)
			if sleepFor < 0 {
				sleepFor = 0
			}
			if sleepFor > policy.MaxDelay {
				sleepFor = policy.MaxDelay
			}
		}
		stats.Delays = append(stats.Delays, sleepFor)
		if err := policy.Sleep(ctx, sleepFor); err != nil {
			return zero, stats, err
		}
		delay *= 2
		if delay > policy.MaxDelay {
			delay = policy.MaxDelay
		}
	}
	return zero, stats, nil
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
