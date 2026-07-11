package cachex

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStoreTTLAndSetNXAndCompareDelete(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	store.Set("session", "token-a", time.Minute)

	if got, ok := store.Get("session"); !ok || got != "token-a" {
		t.Fatalf("Get session = %q/%v, want token-a/true", got, ok)
	}
	if ok := store.SetNX("session", "token-b", time.Minute); ok {
		t.Fatal("SetNX succeeded for existing key")
	}
	if ok := store.CompareAndDelete("session", "token-b"); ok {
		t.Fatal("CompareAndDelete succeeded with wrong token")
	}
	if got, ok := store.Get("session"); !ok || got != "token-a" {
		t.Fatalf("value after wrong delete = %q/%v, want token-a/true", got, ok)
	}
	if ok := store.CompareAndDelete("session", "token-a"); !ok {
		t.Fatal("CompareAndDelete failed for holder token")
	}
	if _, ok := store.Get("session"); ok {
		t.Fatal("key still exists after holder delete")
	}

	store.Set("short", "value", time.Second)
	store.Advance(2 * time.Second)
	if _, ok := store.Get("short"); ok {
		t.Fatal("expired key still exists")
	}
}

func TestStoreExpiresWithRealTime(t *testing.T) {
	store := NewStore(time.Now())
	store.Set("short", "value", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	if _, ok := store.Get("short"); ok {
		t.Fatal("key still exists after real time TTL elapsed")
	}

	lock := NewLockManager(store)
	if ok := lock.Acquire("resource", "token-a", 10*time.Millisecond); !ok {
		t.Fatal("token-a failed to acquire lock")
	}
	time.Sleep(20 * time.Millisecond)
	if ok := lock.Acquire("resource", "token-b", time.Second); !ok {
		t.Fatal("token-b failed to acquire lock after real time TTL elapsed")
	}
}

func TestCacheAsideMissThenHit(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	var loads int
	loader := func(context.Context, string) (string, bool, error) {
		loads++
		return "article-1", true, nil
	}

	got, ok, err := CacheAside(context.Background(), store, "article:1", time.Minute, loader)
	if err != nil || !ok || got != "article-1" {
		t.Fatalf("first CacheAside = %q/%v/%v", got, ok, err)
	}
	got, ok, err = CacheAside(context.Background(), store, "article:1", time.Minute, loader)
	if err != nil || !ok || got != "article-1" {
		t.Fatalf("second CacheAside = %q/%v/%v", got, ok, err)
	}
	if loads != 1 {
		t.Fatalf("loader calls = %d, want 1", loads)
	}
}

func TestReadThroughAndWriteThrough(t *testing.T) {
	source := map[string]string{"user:1": "Ada"}
	readThrough := NewReadThrough(NewStore(time.Unix(100, 0)), time.Minute, func(ctx context.Context, key string) (string, bool, error) {
		value, ok := source[key]
		return value, ok, nil
	})

	got, ok, err := readThrough.Get(context.Background(), "user:1")
	if err != nil || !ok || got != "Ada" {
		t.Fatalf("ReadThrough.Get = %q/%v/%v", got, ok, err)
	}

	writeThrough := NewWriteThrough(NewStore(time.Unix(100, 0)), time.Minute, func(ctx context.Context, key, value string) error {
		source[key] = value
		return nil
	})
	if err := writeThrough.Set(context.Background(), "user:2", "Grace"); err != nil {
		t.Fatalf("WriteThrough.Set returned error: %v", err)
	}
	if source["user:2"] != "Grace" {
		t.Fatalf("source user:2 = %q, want Grace", source["user:2"])
	}
	if got, ok := writeThrough.Store().Get("user:2"); !ok || got != "Grace" {
		t.Fatalf("cache user:2 = %q/%v, want Grace/true", got, ok)
	}
}

func TestNegativeCacheAvoidsRepeatedPenetration(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	var loads int
	loader := func(context.Context, string) (string, bool, error) {
		loads++
		return "", false, nil
	}

	for i := 0; i < 2; i++ {
		got, ok, err := CacheAside(context.Background(), store, "missing", time.Minute, loader, WithNegativeTTL(time.Minute))
		if err != nil || ok || got != "" {
			t.Fatalf("CacheAside missing = %q/%v/%v", got, ok, err)
		}
	}
	if loads != 1 {
		t.Fatalf("loader calls = %d, want 1", loads)
	}
}

func TestJitterTTLSpreadsWithinBounds(t *testing.T) {
	base := 100 * time.Second
	seen := map[time.Duration]bool{}
	for _, key := range []string{"a", "b", "c", "d", "e"} {
		ttl := JitterTTL(base, 10, key)
		if ttl < 90*time.Second || ttl > 110*time.Second {
			t.Fatalf("ttl for %s = %s, want within ±10%%", key, ttl)
		}
		seen[ttl] = true
	}
	if len(seen) < 2 {
		t.Fatalf("jitter produced no spread: %v", seen)
	}
}

func TestJitterTTLNeverReturnsNonPositiveForPositiveBase(t *testing.T) {
	if got := JitterTTL(0, 10, "zero"); got != 0 {
		t.Fatalf("zero base ttl = %s, want 0", got)
	}
	if got := JitterTTL(-time.Second, 10, "negative"); got != 0 {
		t.Fatalf("negative base ttl = %s, want 0", got)
	}
	for _, key := range []string{"key0", "key1", "key2", "key3"} {
		if got := JitterTTL(time.Second, 200, key); got <= 0 {
			t.Fatalf("jitter ttl for %s = %s, want positive", key, got)
		}
	}
}

func TestSingleFlightLoadsOnce(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	var calls int32
	loader := func(context.Context, string) (string, bool, error) {
		atomic.AddInt32(&calls, 1)
		time.Sleep(10 * time.Millisecond)
		return "shared", true, nil
	}
	group := NewSingleFlightCache(store, time.Minute, loader)

	var wg sync.WaitGroup
	results := make(chan string, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, ok, err := group.Get(context.Background(), "hot-key")
			if err != nil || !ok {
				t.Errorf("Get returned %q/%v/%v", got, ok, err)
				return
			}
			results <- got
		}()
	}
	wg.Wait()
	close(results)

	for got := range results {
		if got != "shared" {
			t.Fatalf("result = %q, want shared", got)
		}
	}
	if calls != 1 {
		t.Fatalf("loader calls = %d, want 1", calls)
	}
}

func TestSingleFlightWaiterRespectsContextCancellation(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	started := make(chan struct{})
	release := make(chan struct{})
	loader := func(context.Context, string) (string, bool, error) {
		close(started)
		<-release
		return "shared", true, nil
	}
	group := NewSingleFlightCache(store, time.Minute, loader)

	leaderDone := make(chan error, 1)
	go func() {
		_, _, err := group.Get(context.Background(), "hot-key")
		leaderDone <- err
	}()
	<-started

	ctx, cancel := context.WithCancel(context.Background())
	waiterDone := make(chan error, 1)
	go func() {
		_, _, err := group.Get(ctx, "hot-key")
		waiterDone <- err
	}()
	cancel()

	select {
	case err := <-waiterDone:
		if err != context.Canceled {
			t.Fatalf("waiter err = %v, want context.Canceled", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("waiter did not return after context cancellation")
	}

	close(release)
	if err := <-leaderDone; err != nil {
		t.Fatalf("leader err = %v", err)
	}
}

func TestSingleFlightCleansUpAfterLoaderPanic(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	var calls int32
	group := NewSingleFlightCache(store, time.Minute, func(context.Context, string) (string, bool, error) {
		if atomic.AddInt32(&calls, 1) == 1 {
			panic("boom")
		}
		return "recovered", true, nil
	})

	func() {
		defer func() {
			if recovered := recover(); recovered == nil {
				t.Fatal("first Get did not panic")
			}
		}()
		_, _, _ = group.Get(context.Background(), "hot-key")
	}()

	done := make(chan struct{})
	var got string
	var ok bool
	var err error
	go func() {
		got, ok, err = group.Get(context.Background(), "hot-key")
		close(done)
	}()
	select {
	case <-done:
		if err != nil || !ok || got != "recovered" {
			t.Fatalf("second Get = %q/%v/%v, want recovered/true/nil", got, ok, err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("second Get blocked after loader panic")
	}
}

func TestLockTokenTTLAndRelease(t *testing.T) {
	store := NewStore(time.Unix(100, 0))
	lock := NewLockManager(store)

	if ok := lock.Acquire("article:1", "token-a", time.Second); !ok {
		t.Fatal("token-a failed to acquire lock")
	}
	if ok := lock.Acquire("article:1", "token-b", time.Second); ok {
		t.Fatal("token-b acquired held lock")
	}
	if ok := lock.Release("article:1", "token-b"); ok {
		t.Fatal("token-b released token-a lock")
	}
	if ok := lock.Release("article:1", "token-a"); !ok {
		t.Fatal("token-a failed to release lock")
	}

	if ok := lock.Acquire("article:1", "token-a", time.Second); !ok {
		t.Fatal("token-a failed to reacquire lock")
	}
	store.Advance(2 * time.Second)
	if ok := lock.Acquire("article:1", "token-b", time.Second); !ok {
		t.Fatal("token-b failed to acquire expired lock")
	}
}

func TestLockRejectsInvalidInputs(t *testing.T) {
	lock := NewLockManager(NewStore(time.Unix(100, 0)))
	tests := []struct {
		name     string
		resource string
		token    string
		ttl      time.Duration
	}{
		{name: "empty resource", resource: "", token: "token", ttl: time.Second},
		{name: "blank resource", resource: "   ", token: "token", ttl: time.Second},
		{name: "empty token", resource: "resource", token: "", ttl: time.Second},
		{name: "blank token", resource: "resource", token: "   ", ttl: time.Second},
		{name: "zero ttl", resource: "resource", token: "token", ttl: 0},
		{name: "negative ttl", resource: "resource", token: "token", ttl: -time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := lock.Acquire(tt.resource, tt.token, tt.ttl); ok {
				t.Fatal("Acquire succeeded for invalid input")
			}
		})
	}
	if ok := lock.Release("resource", ""); ok {
		t.Fatal("Release succeeded with empty token")
	}
}
