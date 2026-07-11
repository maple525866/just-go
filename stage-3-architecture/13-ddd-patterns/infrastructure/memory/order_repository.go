// Package memory provides an in-memory OrderRepository adapter.
package memory

import (
	"context"
	"sync"

	"just-go/stage-3-architecture/13-ddd-patterns/application"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*order.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: make(map[string]*order.Order)}
}

func (r *OrderRepository) Get(ctx context.Context, id string) (*order.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	aggregate, ok := r.orders[id]
	if !ok {
		return nil, application.ErrOrderNotFound
	}
	return aggregate.Clone(), nil
}

func (r *OrderRepository) Save(ctx context.Context, aggregate *order.Order, expectedVersion uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if aggregate == nil {
		return application.ErrOrderConflict
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.orders[aggregate.ID()]
	if !exists {
		if expectedVersion != 0 || aggregate.Version() != 1 {
			return application.ErrOrderConflict
		}
		r.orders[aggregate.ID()] = aggregate.Clone()
		return nil
	}
	if stored.Version() != expectedVersion || aggregate.Version() != expectedVersion+1 {
		return application.ErrOrderConflict
	}
	r.orders[aggregate.ID()] = aggregate.Clone()
	return nil
}
