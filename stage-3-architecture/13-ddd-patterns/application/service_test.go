package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

type repositoryMock struct {
	orders  map[string]*order.Order
	getErr  error
	saveErr error
	trace   *[]string
}

func (m *repositoryMock) Get(_ context.Context, id string) (*order.Order, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	aggregate, ok := m.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return aggregate.Clone(), nil
}

func (m *repositoryMock) Save(_ context.Context, aggregate *order.Order, _ uint64) error {
	if m.trace != nil {
		*m.trace = append(*m.trace, "save")
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	m.orders[aggregate.ID()] = aggregate.Clone()
	return nil
}

type publisherMock struct {
	events []domain.Event
	err    error
	trace  *[]string
}

func (m *publisherMock) Publish(_ context.Context, events ...domain.Event) error {
	if m.trace != nil {
		*m.trace = append(*m.trace, "publish")
	}
	m.events = append(m.events, events...)
	return m.err
}

func TestServiceCreateAddGetAndConfirm(t *testing.T) {
	t.Parallel()

	trace := []string{}
	repository := &repositoryMock{orders: make(map[string]*order.Order), trace: &trace}
	publisher := &publisherMock{trace: &trace}
	confirmedAt := time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC)
	service := mustService(t, repository, publisher, order.PercentageDiscount{Percent: 10}, func() time.Time { return confirmedAt })

	address := mustAddress(t)
	created, err := service.Create(context.Background(), CreateOrder{OrderID: "order-1", CustomerID: "customer-1", Address: address})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Version() != 1 {
		t.Fatalf("created version = %d", created.Version())
	}

	price := mustMoney(t, 5000, "CNY")
	added, err := service.AddLine(context.Background(), AddOrderLine{
		OrderID: "order-1", LineID: "line-1", ProductID: "product-1", Name: "Go Book", UnitPrice: price, Quantity: 2,
	})
	if err != nil || len(added.Lines()) != 1 {
		t.Fatalf("AddLine() = %#v, %v", added, err)
	}

	loaded, err := service.Get(context.Background(), "order-1")
	if err != nil || loaded.Lines()[0].ProductID() != "product-1" {
		t.Fatalf("Get() = %#v, %v", loaded, err)
	}
	trace = nil
	confirmed, err := service.Confirm(context.Background(), "order-1")
	if err != nil {
		t.Fatalf("Confirm() error = %v", err)
	}
	if confirmed.Status() != order.StatusConfirmed || confirmed.Total().Minor() != 9000 {
		t.Fatalf("confirmed = status %s, total %d", confirmed.Status(), confirmed.Total().Minor())
	}
	if !reflect.DeepEqual(trace, []string{"save", "publish"}) {
		t.Fatalf("confirm trace = %v", trace)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name() != order.OrderConfirmedName {
		t.Fatalf("published events = %#v", publisher.events)
	}
}

func TestServiceFailurePaths(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		service := mustService(t, &repositoryMock{orders: make(map[string]*order.Order)}, &publisherMock{}, order.NoDiscount{}, time.Now)
		if _, err := service.Get(context.Background(), "missing"); !errors.Is(err, ErrOrderNotFound) {
			t.Fatalf("Get() error = %v", err)
		}
	})

	t.Run("save failure does not publish", func(t *testing.T) {
		aggregate := populatedOrder(t)
		publisher := &publisherMock{}
		repository := &repositoryMock{orders: map[string]*order.Order{aggregate.ID(): aggregate}, saveErr: ErrOrderConflict}
		service := mustService(t, repository, publisher, order.NoDiscount{}, time.Now)
		if _, err := service.Confirm(context.Background(), aggregate.ID()); !errors.Is(err, ErrOrderConflict) {
			t.Fatalf("Confirm() error = %v", err)
		}
		if len(publisher.events) != 0 {
			t.Fatalf("published after save failure: %v", publisher.events)
		}
	})

	t.Run("publish failure retains confirmed result", func(t *testing.T) {
		aggregate := populatedOrder(t)
		repository := &repositoryMock{orders: map[string]*order.Order{aggregate.ID(): aggregate}}
		publisher := &publisherMock{err: errors.New("bus unavailable")}
		service := mustService(t, repository, publisher, order.NoDiscount{}, time.Now)
		confirmed, err := service.Confirm(context.Background(), aggregate.ID())
		if !errors.Is(err, ErrEventPublish) || confirmed.Status() != order.StatusConfirmed {
			t.Fatalf("Confirm() = %#v, %v", confirmed, err)
		}
		stored := repository.orders[aggregate.ID()]
		if stored.Status() != order.StatusConfirmed {
			t.Fatalf("stored status = %s", stored.Status())
		}
	})
}

func mustService(t *testing.T, repository OrderRepository, publisher EventPublisher, policy order.DiscountPolicy, now Clock) *Service {
	t.Helper()
	service, err := NewService(repository, publisher, policy, now)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	return service
}

func mustAddress(t *testing.T) order.Address {
	t.Helper()
	address, err := order.NewAddress("Alice", "1 Go Road", "Beijing", "100000", "CN")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	return address
}

func mustMoney(t *testing.T, minor int64, currency string) order.Money {
	t.Helper()
	money, err := order.NewMoney(minor, currency)
	if err != nil {
		t.Fatalf("NewMoney() error = %v", err)
	}
	return money
}

func populatedOrder(t *testing.T) *order.Order {
	t.Helper()
	aggregate, err := order.New("order-1", "customer-1", mustAddress(t))
	if err != nil {
		t.Fatalf("order.New() error = %v", err)
	}
	line, err := order.NewLine("line-1", "product-1", "Go Book", mustMoney(t, 5000, "CNY"), 1)
	if err != nil {
		t.Fatalf("order.NewLine() error = %v", err)
	}
	if err = aggregate.AddLine(line); err != nil {
		t.Fatalf("AddLine() error = %v", err)
	}
	return aggregate
}
