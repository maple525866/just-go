package application

import (
	"context"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

// OrderRepository persists whole aggregates and checks an expected version.
type OrderRepository interface {
	Get(context.Context, string) (*order.Order, error)
	Save(context.Context, *order.Order, uint64) error
}

// EventPublisher dispatches facts after aggregate persistence succeeds.
type EventPublisher interface {
	Publish(context.Context, ...domain.Event) error
}
