// Package events provides synchronous domain-event adapters for the tutorial.
package events

import (
	"context"
	"fmt"
	"sync"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
)

type Handler func(context.Context, domain.Event) error

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

func (b *Bus) Register(eventName string, handler Handler) error {
	if eventName == "" || handler == nil {
		return fmt.Errorf("event name and handler are required")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
	return nil
}

func (b *Bus) Publish(ctx context.Context, domainEvents ...domain.Event) error {
	for _, event := range domainEvents {
		if err := ctx.Err(); err != nil {
			return err
		}
		if event == nil {
			return fmt.Errorf("publish nil event")
		}
		b.mu.RLock()
		handlers := append([]Handler(nil), b.handlers[event.Name()]...)
		b.mu.RUnlock()
		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				return fmt.Errorf("handle %s: %w", event.Name(), err)
			}
		}
	}
	return nil
}
