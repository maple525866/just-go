package inventory

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStoreWatchReceivesCurrentAndUpdatedStock(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates, err := store.Watch(ctx, "book-1")
	if err != nil {
		t.Fatal(err)
	}
	assertStockUpdate(t, updates, Stock{SKU: "book-1", Quantity: 10, Version: 1})

	if _, err := store.Adjust("book-1", -2); err != nil {
		t.Fatal(err)
	}
	assertStockUpdate(t, updates, Stock{SKU: "book-1", Quantity: 8, Version: 2})
}

func TestStoreWatchRejectsInvalidRequests(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name string
		ctx  context.Context
		sku  string
		want error
	}{
		{name: "blank sku", ctx: context.Background(), sku: " ", want: ErrInvalidStock},
		{name: "missing sku", ctx: context.Background(), sku: "missing", want: ErrStockNotFound},
		{name: "canceled context", ctx: canceled, sku: "book-1", want: context.Canceled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := store.Watch(tt.ctx, tt.sku); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestStoreWatchCancellationClosesChannel(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	updates, err := store.Watch(ctx, "book-1")
	if err != nil {
		t.Fatal(err)
	}
	<-updates
	cancel()

	select {
	case _, ok := <-updates:
		if ok {
			t.Fatal("watch channel remained open")
		}
	case <-time.After(time.Second):
		t.Fatal("watch channel did not close")
	}
}

func TestStoreWatchSlowSubscriberReceivesLatestWithoutBlocking(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates, err := store.Watch(ctx, "book-1")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = store.Adjust("book-1", 1)
		_, _ = store.Adjust("book-1", 1)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Adjust blocked on slow subscriber")
	}
	assertStockUpdate(t, updates, Stock{SKU: "book-1", Quantity: 12, Version: 3})
}

func assertStockUpdate(t *testing.T, updates <-chan Stock, want Stock) {
	t.Helper()
	select {
	case got, ok := <-updates:
		if !ok {
			t.Fatal("updates channel closed")
		}
		if got != want {
			t.Fatalf("stock = %#v, want %#v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %#v", want)
	}
}
