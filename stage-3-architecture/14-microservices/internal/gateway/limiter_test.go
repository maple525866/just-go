package gateway

import (
	"testing"
	"time"
)

func TestLimiterEnforcesLimitPerKey(t *testing.T) {
	now := time.Unix(100, 0)
	limiter := NewLimiter(func() time.Time { return now })
	if !limiter.Allow("client-a", 2, time.Minute) {
		t.Fatal("first allowed request was rejected")
	}
	if !limiter.Allow("client-a", 2, time.Minute) {
		t.Fatal("second allowed request was rejected")
	}
	if limiter.Allow("client-a", 2, time.Minute) {
		t.Fatal("request above limit was allowed")
	}
	if !limiter.Allow("client-b", 2, time.Minute) {
		t.Fatal("separate key shared the limit")
	}
}

func TestLimiterResetsAfterWindow(t *testing.T) {
	now := time.Unix(100, 0)
	limiter := NewLimiter(func() time.Time { return now })
	if !limiter.Allow("client-a", 1, time.Minute) {
		t.Fatal("first request rejected")
	}
	if limiter.Allow("client-a", 1, time.Minute) {
		t.Fatal("second request allowed")
	}
	now = now.Add(time.Minute)
	if !limiter.Allow("client-a", 1, time.Minute) {
		t.Fatal("request after reset rejected")
	}
}

func TestLimiterAppliesDynamicParameters(t *testing.T) {
	now := time.Unix(100, 0)
	limiter := NewLimiter(func() time.Time { return now })
	if !limiter.Allow("client-a", 2, time.Minute) {
		t.Fatal("first request rejected")
	}
	if limiter.Allow("client-a", 1, time.Minute) {
		t.Fatal("lower dynamic limit was not applied")
	}
	if !limiter.Allow("client-a", 1, 2*time.Minute) {
		t.Fatal("new dynamic window did not reset bucket")
	}
}

func TestLimiterRejectsInvalidInputs(t *testing.T) {
	limiter := NewLimiter(nil)
	for _, test := range []struct {
		key    string
		limit  int
		window time.Duration
	}{
		{key: "", limit: 1, window: time.Second},
		{key: "client", limit: 0, window: time.Second},
		{key: "client", limit: 1, window: 0},
	} {
		if limiter.Allow(test.key, test.limit, test.window) {
			t.Fatalf("invalid input was allowed: %#v", test)
		}
	}
}
