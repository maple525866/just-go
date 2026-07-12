package configcenter

import (
	"context"
	"fmt"
	"sync"
)

type MemoryStore struct {
	mu          sync.RWMutex
	current     Snapshot
	watchers    map[uint64]chan Snapshot
	nextWatcher uint64
	closed      bool
}

func NewMemoryStore(initial GatewayConfig) (*MemoryStore, error) {
	initial = normalizeConfig(initial)
	if err := initial.Validate(); err != nil {
		return nil, err
	}
	return &MemoryStore{
		current:  Snapshot{Version: 1, Config: initial},
		watchers: make(map[uint64]chan Snapshot),
	}, nil
}

func (s *MemoryStore) Current() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return Snapshot{}, ErrClosed
	}
	return s.current, nil
}

func (s *MemoryStore) Update(config GatewayConfig) (Snapshot, error) {
	config = normalizeConfig(config)
	if err := config.Validate(); err != nil {
		return Snapshot{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return Snapshot{}, ErrClosed
	}
	s.current = Snapshot{Version: s.current.Version + 1, Config: config}
	s.publishLocked(s.current)
	return s.current, nil
}

func (s *MemoryStore) Watch(ctx context.Context) (<-chan Snapshot, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidConfig)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, ErrClosed
	}
	s.nextWatcher++
	id := s.nextWatcher
	ch := make(chan Snapshot, 1)
	s.watchers[id] = ch
	ch <- s.current
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		if current, exists := s.watchers[id]; exists && current == ch {
			delete(s.watchers, id)
			close(ch)
		}
		s.mu.Unlock()
	}()
	return ch, nil
}

func (s *MemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	for id, ch := range s.watchers {
		close(ch)
		delete(s.watchers, id)
	}
	return nil
}

func (s *MemoryStore) publishLocked(snapshot Snapshot) {
	for _, ch := range s.watchers {
		select {
		case ch <- snapshot:
		default:
			select {
			case <-ch:
			default:
			}
			ch <- snapshot
		}
	}
}
