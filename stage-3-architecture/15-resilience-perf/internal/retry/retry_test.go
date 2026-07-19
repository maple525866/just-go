package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemporary = errors.New("temporary")
var errPermanent = errors.New("permanent")

func TestDoRetriesRetryableErrorsAndReturnsStats(t *testing.T) {
	var slept []time.Duration
	policy := Policy{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    100 * time.Millisecond,
		ShouldRetry: func(err error) bool { return errors.Is(err, errTemporary) },
		Jitter:      func(d time.Duration, attempt int) time.Duration { return d + time.Duration(attempt)*time.Millisecond },
		Sleep: func(ctx context.Context, d time.Duration) error {
			slept = append(slept, d)
			return nil
		},
	}
	calls := 0
	got, stats, err := Do[string](context.Background(), policy, func(ctx context.Context) (string, error) {
		calls++
		if calls < 3 {
			return "", errTemporary
		}
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if got != "ok" || calls != 3 || stats.Attempts != 3 {
		t.Fatalf("got=%q calls=%d stats=%#v", got, calls, stats)
	}
	want := []time.Duration{11 * time.Millisecond, 22 * time.Millisecond}
	if len(slept) != len(want) {
		t.Fatalf("slept = %#v", slept)
	}
	for i := range want {
		if slept[i] != want[i] {
			t.Fatalf("sleep[%d] = %s want %s", i, slept[i], want[i])
		}
	}
}

func TestDoDoesNotRetryPermanentErrors(t *testing.T) {
	policy := Policy{MaxAttempts: 3, BaseDelay: time.Millisecond, ShouldRetry: func(error) bool { return false }}
	calls := 0
	_, stats, err := Do[string](context.Background(), policy, func(ctx context.Context) (string, error) {
		calls++
		return "", errPermanent
	})
	if !errors.Is(err, errPermanent) {
		t.Fatalf("expected permanent error, got %v", err)
	}
	if calls != 1 || stats.Attempts != 1 {
		t.Fatalf("calls=%d stats=%#v", calls, stats)
	}
}

func TestDoClampsJitterToMaxDelayAndZero(t *testing.T) {
	tests := []struct {
		name   string
		jitter func(time.Duration, int) time.Duration
		want   time.Duration
	}{
		{
			name:   "above max",
			jitter: func(d time.Duration, attempt int) time.Duration { return d + 100*time.Millisecond },
			want:   100 * time.Millisecond,
		},
		{
			name:   "below zero",
			jitter: func(d time.Duration, attempt int) time.Duration { return -10 * time.Millisecond },
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slept []time.Duration
			policy := Policy{
				MaxAttempts: 2,
				BaseDelay:   100 * time.Millisecond,
				MaxDelay:    100 * time.Millisecond,
				ShouldRetry: func(error) bool { return true },
				Jitter:      tt.jitter,
				Sleep: func(ctx context.Context, d time.Duration) error {
					slept = append(slept, d)
					return nil
				},
			}
			_, _, err := Do[string](context.Background(), policy, func(ctx context.Context) (string, error) {
				return "", errTemporary
			})
			if !errors.Is(err, errTemporary) {
				t.Fatalf("expected temporary error, got %v", err)
			}
			if len(slept) != 1 || slept[0] != tt.want {
				t.Fatalf("slept = %#v, want [%s]", slept, tt.want)
			}
		})
	}
}
