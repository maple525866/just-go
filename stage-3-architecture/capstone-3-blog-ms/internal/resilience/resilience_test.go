package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLimiterAndRetry(t *testing.T) {
	limiter := NewLimiter(1, 0.01)
	if !limiter.Allow() || limiter.Allow() {
		t.Fatal("limiter did not enforce capacity")
	}
	calls := 0
	got, err := Retry(context.Background(), 2, time.Nanosecond, func(error) bool { return true }, func(context.Context) (string, error) {
		calls++
		if calls == 1 {
			return "", errors.New("temporary")
		}
		return "ok", nil
	})
	if err != nil || got != "ok" || calls != 2 {
		t.Fatalf("got=%q calls=%d err=%v", got, calls, err)
	}
}
