package inventory

import (
	"errors"
	"sync"
	"testing"
)

func TestStoreAdjustIncrementsVersion(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}

	got, err := store.Adjust("book-1", -2)
	if err != nil {
		t.Fatal(err)
	}
	want := Stock{SKU: "book-1", Quantity: 8, Version: 2}
	if got != want {
		t.Fatalf("stock = %#v, want %#v", got, want)
	}
	stored, err := store.Get("book-1")
	if err != nil {
		t.Fatal(err)
	}
	if stored != want {
		t.Fatalf("stored stock = %#v, want %#v", stored, want)
	}
}

func TestNewStoreRejectsInvalidInitialStock(t *testing.T) {
	tests := []map[string]int64{
		{"": 1},
		{" ": 1},
		{"book-1": -1},
	}
	for _, initial := range tests {
		if _, err := NewStore(initial); !errors.Is(err, ErrInvalidStock) {
			t.Fatalf("NewStore(%#v) error = %v", initial, err)
		}
	}
}

func TestStoreRejectsInvalidReadsAndAdjustments(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 2})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		sku   string
		delta int64
		want  error
	}{
		{name: "blank sku", sku: " ", delta: 1, want: ErrInvalidStock},
		{name: "zero delta", sku: "book-1", delta: 0, want: ErrInvalidStock},
		{name: "missing sku", sku: "missing", delta: 1, want: ErrStockNotFound},
		{name: "negative result", sku: "book-1", delta: -3, want: ErrInvalidStock},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := store.Adjust(tt.sku, tt.delta); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
	if _, err := store.Get(" "); !errors.Is(err, ErrInvalidStock) {
		t.Fatalf("blank Get error = %v", err)
	}
	if _, err := store.Get("missing"); !errors.Is(err, ErrStockNotFound) {
		t.Fatalf("missing Get error = %v", err)
	}

	got, err := store.Get("book-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Quantity != 2 || got.Version != 1 {
		t.Fatalf("failed adjustment mutated stock: %#v", got)
	}
}

func TestStoreSupportsConcurrentAdjustments(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 100})
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := store.Adjust("book-1", 1); err != nil {
				t.Errorf("Adjust: %v", err)
			}
		}()
	}
	wg.Wait()

	got, err := store.Get("book-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Quantity != 132 || got.Version != 33 {
		t.Fatalf("stock = %#v", got)
	}
}
