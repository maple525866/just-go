package discovery

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryRegistryResolveIsDeterministic(t *testing.T) {
	registry := NewMemoryRegistry()
	defer registry.Close()
	for _, instance := range []Instance{
		{Service: "product", ID: "b", Address: "127.0.0.1:5002"},
		{Service: "product", ID: "a", Address: "127.0.0.1:5001"},
	} {
		if _, err := registry.Register(instance); err != nil {
			t.Fatal(err)
		}
	}

	got, err := registry.Resolve("product")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "a" || got.Address != "127.0.0.1:5001" {
		t.Fatalf("resolved = %#v", got)
	}
}

func TestMemoryRegistryValidatesAndRejectsDuplicates(t *testing.T) {
	registry := NewMemoryRegistry()
	defer registry.Close()
	tests := []Instance{
		{ID: "one", Address: "127.0.0.1:5001"},
		{Service: "product", Address: "127.0.0.1:5001"},
		{Service: "product", ID: "one"},
		{Service: "product", ID: "one", Address: "not-an-address"},
	}
	for _, instance := range tests {
		if _, err := registry.Register(instance); !errors.Is(err, ErrInvalidInstance) {
			t.Fatalf("Register(%#v) error = %v", instance, err)
		}
	}
	valid := Instance{Service: "product", ID: "one", Address: "127.0.0.1:5001"}
	if _, err := registry.Register(valid); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Register(valid); !errors.Is(err, ErrDuplicateInstance) {
		t.Fatalf("duplicate error = %v", err)
	}
}

func TestMemoryRegistryDeregisterIsIdempotent(t *testing.T) {
	registry := NewMemoryRegistry()
	defer registry.Close()
	deregister, err := registry.Register(Instance{Service: "product", ID: "one", Address: "127.0.0.1:5001"})
	if err != nil {
		t.Fatal(err)
	}
	if err := deregister(); err != nil {
		t.Fatal(err)
	}
	if err := deregister(); err != nil {
		t.Fatalf("second deregister: %v", err)
	}
	if _, err := registry.Resolve("product"); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("Resolve error = %v", err)
	}
}

func TestMemoryRegistryWatchPublishesImmutableLatestSnapshots(t *testing.T) {
	registry := NewMemoryRegistry()
	defer registry.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	snapshots, err := registry.Watch(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}
	assertInstances(t, snapshots, nil)

	deregisterA, err := registry.Register(Instance{Service: "product", ID: "a", Address: "127.0.0.1:5001"})
	if err != nil {
		t.Fatal(err)
	}
	defer deregisterA()
	deregisterB, err := registry.Register(Instance{Service: "product", ID: "b", Address: "127.0.0.1:5002"})
	if err != nil {
		t.Fatal(err)
	}
	defer deregisterB()

	latest := assertInstances(t, snapshots, []string{"a", "b"})
	latest[0].ID = "mutated"
	deregisterB()
	got := assertInstances(t, snapshots, []string{"a"})
	if got[0].ID != "a" {
		t.Fatalf("registry snapshot was mutated: %#v", got)
	}
}

func TestMemoryRegistryWatchCancellationAndClose(t *testing.T) {
	registry := NewMemoryRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	snapshots, err := registry.Watch(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}
	<-snapshots
	cancel()
	assertChannelClosed(t, snapshots)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	snapshots2, err := registry.Watch(ctx2, "product")
	if err != nil {
		t.Fatal(err)
	}
	<-snapshots2
	if err := registry.Close(); err != nil {
		t.Fatal(err)
	}
	if err := registry.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	assertChannelClosed(t, snapshots2)
	if _, err := registry.Resolve("product"); !errors.Is(err, ErrClosed) {
		t.Fatalf("Resolve after Close = %v", err)
	}
	if _, err := registry.Register(Instance{Service: "product", ID: "x", Address: "127.0.0.1:5003"}); !errors.Is(err, ErrClosed) {
		t.Fatalf("Register after Close = %v", err)
	}
	if _, err := registry.Watch(context.Background(), "product"); !errors.Is(err, ErrClosed) {
		t.Fatalf("Watch after Close = %v", err)
	}
}

func TestMemoryRegistrySupportsConcurrentAccess(t *testing.T) {
	registry := NewMemoryRegistry()
	defer registry.Close()
	if _, err := registry.Register(Instance{Service: "product", ID: "seed", Address: "127.0.0.1:5000"}); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := range 32 {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			cleanup, err := registry.Register(Instance{
				Service: "product",
				ID:      fmt.Sprintf("instance-%02d", i),
				Address: fmt.Sprintf("127.0.0.1:%d", 5100+i),
			})
			if err != nil {
				t.Errorf("Register: %v", err)
				return
			}
			if err := cleanup(); err != nil {
				t.Errorf("deregister: %v", err)
			}
		}(i)
		go func() {
			defer wg.Done()
			if _, err := registry.Resolve("product"); err != nil {
				t.Errorf("Resolve: %v", err)
			}
		}()
	}
	wg.Wait()
}

func assertInstances(t *testing.T, snapshots <-chan []Instance, wantIDs []string) []Instance {
	t.Helper()
	select {
	case got, ok := <-snapshots:
		if !ok {
			t.Fatal("snapshot channel closed")
		}
		if len(got) != len(wantIDs) {
			t.Fatalf("snapshot = %#v, want IDs %v", got, wantIDs)
		}
		for i, id := range wantIDs {
			if got[i].ID != id {
				t.Fatalf("snapshot = %#v, want IDs %v", got, wantIDs)
			}
		}
		return got
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for snapshot")
		return nil
	}
}

func assertChannelClosed(t *testing.T, snapshots <-chan []Instance) {
	t.Helper()
	select {
	case _, ok := <-snapshots:
		if ok {
			t.Fatal("channel remained open")
		}
	case <-time.After(time.Second):
		t.Fatal("channel did not close")
	}
}
