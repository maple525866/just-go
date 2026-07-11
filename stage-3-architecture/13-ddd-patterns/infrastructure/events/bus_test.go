package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

type testEvent struct{ name string }

func (e testEvent) Name() string        { return e.name }
func (testEvent) OccurredAt() time.Time { return time.Time{} }

func TestBusRoutesEventsAndPropagatesErrors(t *testing.T) {
	t.Parallel()

	bus := NewBus()
	called := 0
	if err := bus.Register("known", func(context.Context, domain.Event) error {
		called++
		return nil
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := bus.Publish(context.Background(), testEvent{name: "unknown"}, testEvent{name: "known"}); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if called != 1 {
		t.Fatalf("handler calls = %d", called)
	}

	want := errors.New("handler failed")
	if err := bus.Register("known", func(context.Context, domain.Event) error { return want }); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := bus.Publish(context.Background(), testEvent{name: "known"}); !errors.Is(err, want) {
		t.Fatalf("Publish handler error = %v", err)
	}
}

func TestInventoryProjectionHandlesConfirmationIdempotently(t *testing.T) {
	t.Parallel()

	aggregate := confirmedOrder(t)
	domainEvents := aggregate.PullEvents()
	projection := NewInventoryProjection()
	bus := NewBus()
	if err := bus.Register(order.OrderConfirmedName, projection.Handle); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := bus.Publish(context.Background(), domainEvents...); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if err := bus.Publish(context.Background(), domainEvents...); err != nil {
		t.Fatalf("second Publish() error = %v", err)
	}
	if got := projection.Reserved("product-1"); got != 2 {
		t.Fatalf("Reserved() = %d", got)
	}
}

func confirmedOrder(t *testing.T) *order.Order {
	t.Helper()
	address, _ := order.NewAddress("Alice", "1 Go Road", "Beijing", "100000", "CN")
	aggregate, _ := order.New("order-1", "customer-1", address)
	price, _ := order.NewMoney(5000, "CNY")
	line, _ := order.NewLine("line-1", "product-1", "Go Book", price, 2)
	if err := aggregate.AddLine(line); err != nil {
		t.Fatalf("AddLine() error = %v", err)
	}
	if err := aggregate.Confirm(mustTotal(t), time.Now()); err != nil {
		t.Fatalf("Confirm() error = %v", err)
	}
	return aggregate
}

func mustTotal(t *testing.T) order.Money {
	t.Helper()
	total, err := order.NewMoney(10000, "CNY")
	if err != nil {
		t.Fatalf("NewMoney() error = %v", err)
	}
	return total
}
