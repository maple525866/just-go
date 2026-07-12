package configcenter

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestMemoryStoreVersionsAndPublishesUpdates(t *testing.T) {
	store, err := NewMemoryStore(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates, err := store.Watch(ctx)
	if err != nil {
		t.Fatal(err)
	}
	initial := assertConfigUpdate(t, updates)
	if initial.Version != 1 || initial.Config.BearerToken != "teaching-token" {
		t.Fatalf("initial = %#v", initial)
	}

	next := initial.Config
	next.RateLimit = 20
	updated, err := store.Update(next)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != 2 || updated.Config.RateLimit != 20 {
		t.Fatalf("updated = %#v", updated)
	}
	observed := assertConfigUpdate(t, updates)
	if observed != updated {
		t.Fatalf("observed = %#v, want %#v", observed, updated)
	}
	current, err := store.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current != updated {
		t.Fatalf("current = %#v, want %#v", current, updated)
	}
}

func TestMemoryStoreRejectsInvalidUpdateWithoutAdvancingVersion(t *testing.T) {
	store, err := NewMemoryStore(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before, _ := store.Current()
	invalid := before.Config
	invalid.RolloutPercent = 101
	if _, err := store.Update(invalid); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("error = %v, want ErrInvalidConfig", err)
	}
	after, _ := store.Current()
	if after != before {
		t.Fatalf("snapshot changed: before=%#v after=%#v", before, after)
	}
	if _, err := NewMemoryStore(invalid); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("invalid initial error = %v", err)
	}
}

func TestMemoryStoreSlowSubscriberReceivesLatest(t *testing.T) {
	store, err := NewMemoryStore(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates, err := store.Watch(ctx)
	if err != nil {
		t.Fatal(err)
	}

	config := validConfig()
	config.RateLimit = 20
	if _, err := store.Update(config); err != nil {
		t.Fatal(err)
	}
	config.RateLimit = 30
	if _, err := store.Update(config); err != nil {
		t.Fatal(err)
	}
	latest := assertConfigUpdate(t, updates)
	if latest.Version != 3 || latest.Config.RateLimit != 30 {
		t.Fatalf("latest = %#v", latest)
	}
}

func TestMemoryStoreCancellationAndClose(t *testing.T) {
	store, err := NewMemoryStore(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	updates, err := store.Watch(ctx)
	if err != nil {
		t.Fatal(err)
	}
	<-updates
	cancel()
	assertConfigChannelClosed(t, updates)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	updates2, err := store.Watch(ctx2)
	if err != nil {
		t.Fatal(err)
	}
	<-updates2
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	assertConfigChannelClosed(t, updates2)
	if _, err := store.Current(); !errors.Is(err, ErrClosed) {
		t.Fatalf("Current after Close = %v", err)
	}
	if _, err := store.Update(validConfig()); !errors.Is(err, ErrClosed) {
		t.Fatalf("Update after Close = %v", err)
	}
	if _, err := store.Watch(context.Background()); !errors.Is(err, ErrClosed) {
		t.Fatalf("Watch after Close = %v", err)
	}
}

func TestMemoryStoreCloseStopsBackgroundWatchers(t *testing.T) {
	store, err := NewMemoryStore(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	updates, err := store.Watch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	<-updates

	done := make(chan struct{})
	go func() {
		store.watchWG.Wait()
		close(done)
	}()
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("background watcher did not stop after Close")
	}
}

func TestMemoryStoreSupportsConcurrentAccess(t *testing.T) {
	store, err := NewMemoryStore(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var wg sync.WaitGroup
	for i := range 32 {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			config := validConfig()
			config.RateLimit = i + 1
			if _, err := store.Update(config); err != nil {
				t.Errorf("Update: %v", err)
			}
		}(i)
		go func() {
			defer wg.Done()
			if _, err := store.Current(); err != nil {
				t.Errorf("Current: %v", err)
			}
		}()
	}
	wg.Wait()
	current, err := store.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.Version != 33 {
		t.Fatalf("version = %d, want 33", current.Version)
	}
}

func assertConfigUpdate(t *testing.T, updates <-chan Snapshot) Snapshot {
	t.Helper()
	select {
	case update, ok := <-updates:
		if !ok {
			t.Fatal("updates channel closed")
		}
		return update
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for config update")
		return Snapshot{}
	}
}

func assertConfigChannelClosed(t *testing.T, updates <-chan Snapshot) {
	t.Helper()
	select {
	case _, ok := <-updates:
		if ok {
			t.Fatal("channel remained open")
		}
	case <-time.After(time.Second):
		t.Fatal("channel did not close")
	}
}
