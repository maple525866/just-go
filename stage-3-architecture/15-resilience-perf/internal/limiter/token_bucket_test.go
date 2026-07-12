package limiter

import (
	"math"
	"testing"
	"time"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time          { return c.now }
func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

func TestTokenBucketAllowsBurstThenRejectsWithRetryAfter(t *testing.T) {
	clock := &fakeClock{now: time.Unix(100, 0)}
	bucket, err := NewTokenBucket(Config{Capacity: 2, RefillPerSecond: 1}, clock)
	if err != nil {
		t.Fatalf("NewTokenBucket: %v", err)
	}

	if decision := bucket.Allow(); !decision.Allowed || decision.Remaining != 1 {
		t.Fatalf("first decision = %#v", decision)
	}
	if decision := bucket.Allow(); !decision.Allowed || decision.Remaining != 0 {
		t.Fatalf("second decision = %#v", decision)
	}
	decision := bucket.Allow()
	if decision.Allowed {
		t.Fatalf("third decision should reject: %#v", decision)
	}
	if decision.RetryAfter != time.Second {
		t.Fatalf("retry after = %s", decision.RetryAfter)
	}
}

func TestTokenBucketRefillsOverInjectedTime(t *testing.T) {
	clock := &fakeClock{now: time.Unix(100, 0)}
	bucket, err := NewTokenBucket(Config{Capacity: 3, RefillPerSecond: 2}, clock)
	if err != nil {
		t.Fatalf("NewTokenBucket: %v", err)
	}

	for i := 0; i < 3; i++ {
		if decision := bucket.Allow(); !decision.Allowed {
			t.Fatalf("initial token %d rejected: %#v", i, decision)
		}
	}
	if decision := bucket.Allow(); decision.Allowed {
		t.Fatalf("empty bucket allowed: %#v", decision)
	}

	clock.Advance(500 * time.Millisecond)
	if decision := bucket.Allow(); !decision.Allowed || decision.Remaining != 0 {
		t.Fatalf("half-second refill decision = %#v", decision)
	}

	clock.Advance(10 * time.Second)
	if decision := bucket.Allow(); !decision.Allowed || decision.Remaining != 2 {
		t.Fatalf("bucket should cap at capacity: %#v", decision)
	}
}

func TestTokenBucketRejectsInvalidConfig(t *testing.T) {
	if _, err := NewTokenBucket(Config{Capacity: 0, RefillPerSecond: 1}, realClock{}); err == nil {
		t.Fatal("expected invalid capacity error")
	}
	if _, err := NewTokenBucket(Config{Capacity: 1, RefillPerSecond: 0}, realClock{}); err == nil {
		t.Fatal("expected invalid refill rate error")
	}
	if _, err := NewTokenBucket(Config{Capacity: 1, RefillPerSecond: math.NaN()}, realClock{}); err == nil {
		t.Fatal("expected NaN refill rate error")
	}
	if _, err := NewTokenBucket(Config{Capacity: 1, RefillPerSecond: math.Inf(1)}, realClock{}); err == nil {
		t.Fatal("expected infinite refill rate error")
	}
}
