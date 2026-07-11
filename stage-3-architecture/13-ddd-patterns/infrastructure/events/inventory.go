package events

import (
	"context"
	"fmt"
	"sync"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
)

// InventoryProjection is deliberately outside the Order aggregate. It models
// a separate consumer reacting to an immutable domain fact.
type InventoryProjection struct {
	mu        sync.RWMutex
	reserved  map[string]int
	processed map[string]struct{}
}

func NewInventoryProjection() *InventoryProjection {
	return &InventoryProjection{
		reserved:  make(map[string]int),
		processed: make(map[string]struct{}),
	}
}

func (p *InventoryProjection) Handle(_ context.Context, event domain.Event) error {
	confirmed, ok := event.(order.OrderConfirmed)
	if !ok {
		return fmt.Errorf("expected order.OrderConfirmed, got %T", event)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.processed[confirmed.OrderID()]; exists {
		return nil
	}
	for _, line := range confirmed.Lines() {
		p.reserved[line.ProductID()] += line.Quantity()
	}
	p.processed[confirmed.OrderID()] = struct{}{}
	return nil
}

func (p *InventoryProjection) Reserved(productID string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.reserved[productID]
}
