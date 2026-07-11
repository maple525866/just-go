package memory

import (
	"context"
	"errors"
	"testing"

	"just-go/stage-3-architecture/13-ddd-patterns/application"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

func TestOrderRepositorySaveGetAndIsolation(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()
	aggregate := newOrder(t)
	if err := repository.Save(context.Background(), aggregate, 0); err != nil {
		t.Fatalf("Save(create) error = %v", err)
	}

	line := newLine(t, "line-1", "product-1")
	if err := aggregate.AddLine(line); err != nil {
		t.Fatalf("AddLine() error = %v", err)
	}
	loaded, err := repository.Get(context.Background(), aggregate.ID())
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(loaded.Lines()) != 0 {
		t.Fatal("mutating original changed stored aggregate")
	}

	if err = loaded.AddLine(line); err != nil {
		t.Fatalf("loaded AddLine() error = %v", err)
	}
	if err = repository.Save(context.Background(), loaded, 1); err != nil {
		t.Fatalf("Save(update) error = %v", err)
	}
	loadedLines := loaded.Lines()
	loadedLines[0] = order.Line{}
	reloaded, err := repository.Get(context.Background(), aggregate.ID())
	if err != nil || len(reloaded.Lines()) != 1 || reloaded.Lines()[0].ID() != "line-1" {
		t.Fatalf("reloaded = %#v, %v", reloaded, err)
	}
}

func TestOrderRepositoryConflictsAndNotFound(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()
	if _, err := repository.Get(context.Background(), "missing"); !errors.Is(err, application.ErrOrderNotFound) {
		t.Fatalf("Get missing error = %v", err)
	}
	aggregate := newOrder(t)
	if err := repository.Save(context.Background(), aggregate, 0); err != nil {
		t.Fatalf("Save(create) error = %v", err)
	}
	first, _ := repository.Get(context.Background(), aggregate.ID())
	stale, _ := repository.Get(context.Background(), aggregate.ID())
	if err := first.AddLine(newLine(t, "line-1", "product-1")); err != nil {
		t.Fatalf("first AddLine() error = %v", err)
	}
	if err := repository.Save(context.Background(), first, 1); err != nil {
		t.Fatalf("first Save() error = %v", err)
	}
	if err := stale.AddLine(newLine(t, "line-2", "product-2")); err != nil {
		t.Fatalf("stale AddLine() error = %v", err)
	}
	if err := repository.Save(context.Background(), stale, 1); !errors.Is(err, application.ErrOrderConflict) {
		t.Fatalf("stale Save() error = %v", err)
	}
	if err := repository.Save(context.Background(), aggregate, 0); !errors.Is(err, application.ErrOrderConflict) {
		t.Fatalf("duplicate Save() error = %v", err)
	}
}

func newOrder(t *testing.T) *order.Order {
	t.Helper()
	address, err := order.NewAddress("Alice", "1 Go Road", "Beijing", "100000", "CN")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	aggregate, err := order.New("order-1", "customer-1", address)
	if err != nil {
		t.Fatalf("order.New() error = %v", err)
	}
	return aggregate
}

func newLine(t *testing.T, lineID, productID string) order.Line {
	t.Helper()
	price, err := order.NewMoney(5000, "CNY")
	if err != nil {
		t.Fatalf("NewMoney() error = %v", err)
	}
	line, err := order.NewLine(lineID, productID, "Go Book", price, 1)
	if err != nil {
		t.Fatalf("NewLine() error = %v", err)
	}
	return line
}
