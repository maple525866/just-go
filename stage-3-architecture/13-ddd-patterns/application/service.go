package application

import (
	"context"
	"fmt"
	"time"

	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

type Clock func() time.Time

type Service struct {
	repository OrderRepository
	publisher  EventPublisher
	pricing    order.PricingService
	policy     order.DiscountPolicy
	now        Clock
}

func NewService(repository OrderRepository, publisher EventPublisher, policy order.DiscountPolicy, now Clock) (*Service, error) {
	if repository == nil || publisher == nil || policy == nil || now == nil {
		return nil, ErrDependency
	}
	return &Service{repository: repository, publisher: publisher, policy: policy, now: now}, nil
}

type CreateOrder struct {
	OrderID    string
	CustomerID string
	Address    order.Address
}

func (s *Service) Create(ctx context.Context, command CreateOrder) (*order.Order, error) {
	aggregate, err := order.New(command.OrderID, command.CustomerID, command.Address)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, aggregate, 0); err != nil {
		return nil, fmt.Errorf("save new order: %w", err)
	}
	return aggregate.Clone(), nil
}

type AddOrderLine struct {
	OrderID   string
	LineID    string
	ProductID string
	Name      string
	UnitPrice order.Money
	Quantity  int
}

func (s *Service) AddLine(ctx context.Context, command AddOrderLine) (*order.Order, error) {
	aggregate, err := s.repository.Get(ctx, command.OrderID)
	if err != nil {
		return nil, fmt.Errorf("load order: %w", err)
	}
	expectedVersion := aggregate.Version()
	line, err := order.NewLine(command.LineID, command.ProductID, command.Name, command.UnitPrice, command.Quantity)
	if err != nil {
		return nil, err
	}
	if err = aggregate.AddLine(line); err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, aggregate, expectedVersion); err != nil {
		return nil, fmt.Errorf("save order line: %w", err)
	}
	return aggregate.Clone(), nil
}

func (s *Service) Get(ctx context.Context, orderID string) (*order.Order, error) {
	aggregate, err := s.repository.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("load order: %w", err)
	}
	return aggregate.Clone(), nil
}

func (s *Service) Confirm(ctx context.Context, orderID string) (*order.Order, error) {
	aggregate, err := s.repository.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("load order: %w", err)
	}
	expectedVersion := aggregate.Version()
	total, err := s.pricing.Calculate(aggregate.Lines(), s.policy)
	if err != nil {
		return nil, err
	}
	if err = aggregate.Confirm(total, s.now()); err != nil {
		return nil, err
	}
	if err = s.repository.Save(ctx, aggregate, expectedVersion); err != nil {
		return nil, fmt.Errorf("save confirmed order: %w", err)
	}
	events := aggregate.PullEvents()
	if err = s.publisher.Publish(ctx, events...); err != nil {
		return aggregate.Clone(), fmt.Errorf("%w: %v", ErrEventPublish, err)
	}
	return aggregate.Clone(), nil
}
