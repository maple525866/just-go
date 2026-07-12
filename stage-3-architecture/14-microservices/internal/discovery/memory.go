package discovery

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
)

type MemoryRegistry struct {
	mu          sync.RWMutex
	services    map[string]map[string]Instance
	watchers    map[string]map[uint64]chan []Instance
	nextWatcher uint64
	closed      bool
}

func NewMemoryRegistry() *MemoryRegistry {
	return &MemoryRegistry{
		services: make(map[string]map[string]Instance),
		watchers: make(map[string]map[uint64]chan []Instance),
	}
}

func (r *MemoryRegistry) Register(candidate Instance) (func() error, error) {
	instance, err := normalizeInstance(candidate)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil, ErrClosed
	}
	instances := r.services[instance.Service]
	if instances == nil {
		instances = make(map[string]Instance)
		r.services[instance.Service] = instances
	}
	if _, exists := instances[instance.ID]; exists {
		r.mu.Unlock()
		return nil, fmt.Errorf("%w: %s/%s", ErrDuplicateInstance, instance.Service, instance.ID)
	}
	instances[instance.ID] = instance
	r.publishLocked(instance.Service)
	r.mu.Unlock()

	var once sync.Once
	return func() error {
		var deregisterErr error
		once.Do(func() {
			deregisterErr = r.deregister(instance.Service, instance.ID)
		})
		return deregisterErr
	}, nil
}

func (r *MemoryRegistry) Resolve(service string) (Instance, error) {
	service = strings.TrimSpace(service)
	if service == "" {
		return Instance{}, fmt.Errorf("%w: service is required", ErrInvalidInstance)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.closed {
		return Instance{}, ErrClosed
	}
	instances := r.snapshotLocked(service)
	if len(instances) == 0 {
		return Instance{}, fmt.Errorf("%w: %s", ErrUnavailable, service)
	}
	return instances[0], nil
}

func (r *MemoryRegistry) Watch(ctx context.Context, service string) (<-chan []Instance, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidInstance)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	service = strings.TrimSpace(service)
	if service == "" {
		return nil, fmt.Errorf("%w: service is required", ErrInvalidInstance)
	}

	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil, ErrClosed
	}
	r.nextWatcher++
	id := r.nextWatcher
	ch := make(chan []Instance, 1)
	if r.watchers[service] == nil {
		r.watchers[service] = make(map[uint64]chan []Instance)
	}
	r.watchers[service][id] = ch
	ch <- r.snapshotLocked(service)
	r.mu.Unlock()

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		watchers := r.watchers[service]
		if current, exists := watchers[id]; exists && current == ch {
			delete(watchers, id)
			if len(watchers) == 0 {
				delete(r.watchers, service)
			}
			close(ch)
		}
		r.mu.Unlock()
	}()
	return ch, nil
}

func (r *MemoryRegistry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	for _, watchers := range r.watchers {
		for _, ch := range watchers {
			close(ch)
		}
	}
	r.watchers = make(map[string]map[uint64]chan []Instance)
	return nil
}

func (r *MemoryRegistry) deregister(service, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	instances := r.services[service]
	if _, exists := instances[id]; !exists {
		return nil
	}
	delete(instances, id)
	if len(instances) == 0 {
		delete(r.services, service)
	}
	r.publishLocked(service)
	return nil
}

func (r *MemoryRegistry) publishLocked(service string) {
	snapshot := r.snapshotLocked(service)
	for _, ch := range r.watchers[service] {
		copyForWatcher := append([]Instance(nil), snapshot...)
		select {
		case ch <- copyForWatcher:
		default:
			select {
			case <-ch:
			default:
			}
			ch <- copyForWatcher
		}
	}
}

func (r *MemoryRegistry) snapshotLocked(service string) []Instance {
	instances := r.services[service]
	snapshot := make([]Instance, 0, len(instances))
	for _, instance := range instances {
		snapshot = append(snapshot, instance)
	}
	sort.Slice(snapshot, func(i, j int) bool {
		if snapshot[i].ID == snapshot[j].ID {
			return snapshot[i].Address < snapshot[j].Address
		}
		return snapshot[i].ID < snapshot[j].ID
	})
	return snapshot
}

func normalizeInstance(candidate Instance) (Instance, error) {
	instance := Instance{
		Service: strings.TrimSpace(candidate.Service),
		ID:      strings.TrimSpace(candidate.ID),
		Address: strings.TrimSpace(candidate.Address),
	}
	if instance.Service == "" || instance.ID == "" || instance.Address == "" {
		return Instance{}, fmt.Errorf("%w: service, id, and address are required", ErrInvalidInstance)
	}
	if _, _, err := net.SplitHostPort(instance.Address); err != nil {
		return Instance{}, fmt.Errorf("%w: address: %v", ErrInvalidInstance, err)
	}
	return instance, nil
}
