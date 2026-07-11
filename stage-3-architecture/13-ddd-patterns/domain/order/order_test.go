package order

import (
	"errors"
	"testing"
	"time"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
)

func TestOrderAggregateLifecycleAndEvents(t *testing.T) {
	t.Parallel()

	aggregate, err := New("order-1", "customer-1", mustAddress(t))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if aggregate.Status() != StatusDraft || aggregate.Version() != 1 {
		t.Fatalf("new order status/version = %s/%d", aggregate.Status(), aggregate.Version())
	}

	line, err := NewLine("line-1", "product-1", "Go Book", mustMoney(t, 5000, "CNY"), 2)
	if err != nil {
		t.Fatalf("NewLine() error = %v", err)
	}
	if err = aggregate.AddLine(line); err != nil {
		t.Fatalf("AddLine() error = %v", err)
	}
	if !aggregate.Total().Equal(mustMoney(t, 10000, "CNY")) {
		t.Fatalf("draft total = %d", aggregate.Total().Minor())
	}
	if err = aggregate.AddLine(line); !errors.Is(err, domain.ErrDuplicateLine) {
		t.Fatalf("duplicate AddLine() error = %v", err)
	}

	lines := aggregate.Lines()
	lines[0] = Line{}
	if aggregate.Lines()[0].ID() != "line-1" {
		t.Fatal("Lines exposed aggregate storage")
	}

	confirmedAt := time.Date(2026, time.July, 12, 10, 0, 0, 0, time.UTC)
	total := mustMoney(t, 9000, "CNY")
	if err = aggregate.Confirm(total, confirmedAt); err != nil {
		t.Fatalf("Confirm() error = %v", err)
	}
	if aggregate.Status() != StatusConfirmed || !aggregate.Total().Equal(total) || aggregate.Version() != 3 {
		t.Fatalf("confirmed order = status %s, total %#v, version %d", aggregate.Status(), aggregate.Total(), aggregate.Version())
	}
	if err = aggregate.AddLine(line); !errors.Is(err, domain.ErrInvalidState) {
		t.Fatalf("confirmed AddLine() error = %v", err)
	}
	if err = aggregate.RemoveLine("line-1"); !errors.Is(err, domain.ErrInvalidState) {
		t.Fatalf("confirmed RemoveLine() error = %v", err)
	}

	events := aggregate.PullEvents()
	if len(events) != 1 || events[0].Name() != OrderConfirmedName {
		t.Fatalf("events = %#v", events)
	}
	event, ok := events[0].(OrderConfirmed)
	if !ok || event.OrderID() != "order-1" || len(event.Lines()) != 1 || event.Lines()[0].Quantity() != 2 {
		t.Fatalf("OrderConfirmed = %#v", events[0])
	}
	if len(aggregate.PullEvents()) != 0 {
		t.Fatal("PullEvents did not clear pending events")
	}
}

func TestOrderAggregateRejectsInvalidTransitions(t *testing.T) {
	t.Parallel()

	if _, err := New("", "customer", mustAddress(t)); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("New invalid error = %v", err)
	}
	if _, err := New("order", "customer", Address{}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("New empty address error = %v", err)
	}
	if _, err := NewLine("", "product", "name", mustMoney(t, 1, "CNY"), 1); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("NewLine invalid error = %v", err)
	}

	empty, err := New("order-1", "customer-1", mustAddress(t))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err = empty.Confirm(mustMoney(t, 0, "CNY"), time.Now()); !errors.Is(err, domain.ErrInvalidState) {
		t.Fatalf("empty Confirm() error = %v", err)
	}
	if err = empty.RemoveLine("missing"); !errors.Is(err, domain.ErrLineNotFound) {
		t.Fatalf("RemoveLine() error = %v", err)
	}
	if err = empty.AddLine(Line{}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("AddLine zero value error = %v", err)
	}
}

func TestOrderMaintainsDraftTotal(t *testing.T) {
	t.Parallel()

	aggregate, err := New("order-1", "customer-1", mustAddress(t))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	first, _ := NewLine("line-1", "product-1", "One", mustMoney(t, 1000, "CNY"), 2)
	second, _ := NewLine("line-2", "product-2", "Two", mustMoney(t, 500, "CNY"), 1)
	foreign, _ := NewLine("line-3", "product-3", "Three", mustMoney(t, 100, "USD"), 1)
	if err = aggregate.AddLine(first); err != nil {
		t.Fatalf("AddLine(first) error = %v", err)
	}
	if err = aggregate.AddLine(second); err != nil {
		t.Fatalf("AddLine(second) error = %v", err)
	}
	if aggregate.Total().Minor() != 2500 {
		t.Fatalf("total after additions = %d", aggregate.Total().Minor())
	}
	if err = aggregate.AddLine(foreign); !errors.Is(err, domain.ErrCurrencyMismatch) {
		t.Fatalf("AddLine currency error = %v", err)
	}
	if err = aggregate.RemoveLine("line-1"); err != nil {
		t.Fatalf("RemoveLine() error = %v", err)
	}
	if aggregate.Total().Minor() != 500 || len(aggregate.Lines()) != 1 {
		t.Fatalf("after removal: total=%d lines=%d", aggregate.Total().Minor(), len(aggregate.Lines()))
	}
	if err = aggregate.RemoveLine("line-2"); err != nil {
		t.Fatalf("RemoveLine(last) error = %v", err)
	}
	if aggregate.Total().Minor() != 0 || aggregate.Total().Currency() != "CNY" {
		t.Fatalf("empty total = %d %s", aggregate.Total().Minor(), aggregate.Total().Currency())
	}
}

func TestPricingService(t *testing.T) {
	t.Parallel()

	lineOne, _ := NewLine("1", "p1", "One", mustMoney(t, 1000, "CNY"), 2)
	lineTwo, _ := NewLine("2", "p2", "Two", mustMoney(t, 500, "CNY"), 1)
	total, err := (PricingService{}).Calculate([]Line{lineOne, lineTwo}, PercentageDiscount{Percent: 10})
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	if !total.Equal(mustMoney(t, 2250, "CNY")) {
		t.Fatalf("Calculate() = %d", total.Minor())
	}
	if _, err = (PricingService{}).Calculate(nil, NoDiscount{}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("empty Calculate() error = %v", err)
	}
	if _, err = (PricingService{}).Calculate([]Line{lineOne}, PercentageDiscount{Percent: 101}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("invalid discount error = %v", err)
	}
}
