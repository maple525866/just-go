package bulkhead

import (
	"context"
	"errors"
	"testing"
)

func TestBulkheadRejectsWhenFullAndReleasesSlots(t *testing.T) {
	b, err := New(1)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	release, err := b.Acquire(context.Background())
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	if _, err := b.Acquire(context.Background()); !errors.Is(err, ErrFull) {
		t.Fatalf("expected ErrFull, got %v", err)
	}
	release()
	release2, err := b.Acquire(context.Background())
	if err != nil {
		t.Fatalf("acquire after release: %v", err)
	}
	release2()
}

func TestBulkheadHonorsCanceledContext(t *testing.T) {
	b, err := New(1)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	release, err := b.Acquire(context.Background())
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	defer release()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := b.Acquire(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context, got %v", err)
	}
}
