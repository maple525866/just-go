package inventory

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrInvalidStock  = errors.New("invalid stock")
	ErrStockNotFound = errors.New("stock not found")
)

type Stock struct {
	SKU      string
	Quantity int64
	Version  uint64
}

type Store struct {
	mu          sync.RWMutex
	stock       map[string]Stock
	watchers    map[string]map[uint64]chan Stock
	nextWatcher uint64
}

func NewStore(initial map[string]int64) (*Store, error) {
	store := &Store{
		stock:    make(map[string]Stock, len(initial)),
		watchers: make(map[string]map[uint64]chan Stock),
	}
	for rawSKU, quantity := range initial {
		sku := strings.TrimSpace(rawSKU)
		if sku == "" || quantity < 0 {
			return nil, fmt.Errorf("%w: sku and non-negative quantity are required", ErrInvalidStock)
		}
		if _, exists := store.stock[sku]; exists {
			return nil, fmt.Errorf("%w: duplicate sku %q", ErrInvalidStock, sku)
		}
		store.stock[sku] = Stock{SKU: sku, Quantity: quantity, Version: 1}
	}
	return store, nil
}

func (s *Store) Get(sku string) (Stock, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return Stock{}, fmt.Errorf("%w: sku is required", ErrInvalidStock)
	}

	s.mu.RLock()
	stock, ok := s.stock[sku]
	s.mu.RUnlock()
	if !ok {
		return Stock{}, fmt.Errorf("%w: %s", ErrStockNotFound, sku)
	}
	return stock, nil
}

func (s *Store) Adjust(sku string, delta int64) (Stock, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" || delta == 0 {
		return Stock{}, fmt.Errorf("%w: sku and non-zero delta are required", ErrInvalidStock)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	stock, ok := s.stock[sku]
	if !ok {
		return Stock{}, fmt.Errorf("%w: %s", ErrStockNotFound, sku)
	}
	if stock.Quantity+delta < 0 {
		return Stock{}, fmt.Errorf("%w: quantity cannot be negative", ErrInvalidStock)
	}
	stock.Quantity += delta
	stock.Version++
	s.stock[sku] = stock
	s.publishLocked(stock)
	return stock, nil
}

func (s *Store) Watch(ctx context.Context, sku string) (<-chan Stock, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidStock)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return nil, fmt.Errorf("%w: sku is required", ErrInvalidStock)
	}

	s.mu.Lock()
	stock, ok := s.stock[sku]
	if !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf("%w: %s", ErrStockNotFound, sku)
	}
	s.nextWatcher++
	id := s.nextWatcher
	ch := make(chan Stock, 1)
	if s.watchers[sku] == nil {
		s.watchers[sku] = make(map[uint64]chan Stock)
	}
	s.watchers[sku][id] = ch
	ch <- stock
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		watchers := s.watchers[sku]
		if current, exists := watchers[id]; exists && current == ch {
			delete(watchers, id)
			if len(watchers) == 0 {
				delete(s.watchers, sku)
			}
			close(ch)
		}
		s.mu.Unlock()
	}()
	return ch, nil
}

func (s *Store) publishLocked(stock Stock) {
	for _, ch := range s.watchers[stock.SKU] {
		select {
		case ch <- stock:
		default:
			select {
			case <-ch:
			default:
			}
			ch <- stock
		}
	}
}
