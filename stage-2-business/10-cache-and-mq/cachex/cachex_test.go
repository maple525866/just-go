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
