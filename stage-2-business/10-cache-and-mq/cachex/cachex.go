package cachex

import (
	"context"
	"hash/fnv"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
	negative  bool
}

// Store is a Redis-like in-memory key-value store with TTL.
type Store struct {
	mu    sync.Mutex
	now   time.Time
	items map[string]entry
}

// NewStore creates a store whose clock can be advanced in tests.
func NewStore(now time.Time) *Store {
	return &Store{now: now, items: map[string]entry{}}
}

// Advance moves the store clock forward.
func (s *Store) Advance(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.now = s.now.Add(d)
}

// Get returns a non-negative cached value.
func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.getLocked(key)
	if !ok || item.negative {
		return "", false
	}
	return item.value, true
}

func (s *Store) getEntry(key string) (entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getLocked(key)
}

func (s *Store) getLocked(key string) (entry, bool) {
	item, ok := s.items[key]
	if !ok {
		return entry{}, false
	}
	if !item.expiresAt.IsZero() && !s.now.Before(item.expiresAt) {
		delete(s.items, key)
		return entry{}, false
	}
	return item, true
}

// Set writes a value with TTL.
func (s *Store) Set(key, value string, ttl time.Duration) {
	s.set(key, value, ttl, false)
}

func (s *Store) setNegative(key string, ttl time.Duration) {
	s.set(key, "", ttl, true)
}

func (s *Store) set(key, value string, ttl time.Duration, negative bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = s.now.Add(ttl)
	}
	s.items[key] = entry{value: value, expiresAt: expiresAt, negative: negative}
}

// Delete removes a key.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

// SetNX sets a key only when it does not exist.
func (s *Store) SetNX(key, value string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.getLocked(key); ok {
		return false
	}
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = s.now.Add(ttl)
	}
	s.items[key] = entry{value: value, expiresAt: expiresAt}
	return true
}

// CompareAndDelete removes a key only when the stored token matches.
func (s *Store) CompareAndDelete(key, token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.getLocked(key)
	if !ok || item.value != token {
		return false
	}
	delete(s.items, key)
	return true
}

// Loader fetches data from the source of truth.
type Loader func(context.Context, string) (string, bool, error)

type cacheAsideOptions struct {
	negativeTTL time.Duration
}

// Option adjusts cache behavior.
type Option func(*cacheAsideOptions)

// WithNegativeTTL caches missing source values for ttl.
func WithNegativeTTL(ttl time.Duration) Option {
	return func(options *cacheAsideOptions) { options.negativeTTL = ttl }
}

// CacheAside implements the cache-aside pattern.
func CacheAside(ctx context.Context, store *Store, key string, ttl time.Duration, loader Loader, options ...Option) (string, bool, error) {
	if item, ok := store.getEntry(key); ok {
		return item.value, !item.negative, nil
	}
	config := cacheAsideOptions{}
	for _, option := range options {
		option(&config)
	}
	value, ok, err := loader(ctx, key)
	if err != nil {
		return "", false, err
	}
	if !ok {
		if config.negativeTTL > 0 {
			store.setNegative(key, config.negativeTTL)
		}
		return "", false, nil
	}
	store.Set(key, value, ttl)
	return value, true, nil
}

// ReadThrough owns its loader and fills cache on miss.
type ReadThrough struct {
	store  *Store
	ttl    time.Duration
	loader Loader
}

func NewReadThrough(store *Store, ttl time.Duration, loader Loader) *ReadThrough {
	return &ReadThrough{store: store, ttl: ttl, loader: loader}
}

func (c *ReadThrough) Get(ctx context.Context, key string) (string, bool, error) {
	return CacheAside(ctx, c.store, key, c.ttl, c.loader)
}

// WriteThrough writes source and cache together.
type WriteThrough struct {
	store *Store
	ttl   time.Duration
	write func(context.Context, string, string) error
}

func NewWriteThrough(store *Store, ttl time.Duration, write func(context.Context, string, string) error) *WriteThrough {
	return &WriteThrough{store: store, ttl: ttl, write: write}
}

func (c *WriteThrough) Set(ctx context.Context, key, value string) error {
	if err := c.write(ctx, key, value); err != nil {
		return err
	}
	c.store.Set(key, value, c.ttl)
	return nil
}

func (c *WriteThrough) Store() *Store { return c.store }

// JitterTTL returns a deterministic TTL spread by percent for a key.
func JitterTTL(base time.Duration, percent int, key string) time.Duration {
	if percent <= 0 {
		return base
	}
	rangeNanos := int64(base) * int64(percent) / 100
	if rangeNanos == 0 {
		return base
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	offset := int64(h.Sum32())%(2*rangeNanos+1) - rangeNanos
	return base + time.Duration(offset)
}

// SingleFlightCache coalesces concurrent cache misses for the same key.
type SingleFlightCache struct {
	store  *Store
	ttl    time.Duration
	loader Loader
	mu     sync.Mutex
	calls  map[string]*call
}

type call struct {
	wg    sync.WaitGroup
	value string
	ok    bool
	err   error
}

func NewSingleFlightCache(store *Store, ttl time.Duration, loader Loader) *SingleFlightCache {
	return &SingleFlightCache{store: store, ttl: ttl, loader: loader, calls: map[string]*call{}}
}

func (c *SingleFlightCache) Get(ctx context.Context, key string) (string, bool, error) {
	if value, ok := c.store.Get(key); ok {
		return value, true, nil
	}

	c.mu.Lock()
	if existing, ok := c.calls[key]; ok {
		c.mu.Unlock()
		existing.wg.Wait()
		return existing.value, existing.ok, existing.err
	}
	current := &call{}
	current.wg.Add(1)
	c.calls[key] = current
	c.mu.Unlock()

	current.value, current.ok, current.err = CacheAside(ctx, c.store, key, c.ttl, c.loader)
	current.wg.Done()

	c.mu.Lock()
	delete(c.calls, key)
	c.mu.Unlock()
	return current.value, current.ok, current.err
}

// LockManager implements a Redis-style token lock.
type LockManager struct {
	store *Store
}

func NewLockManager(store *Store) *LockManager {
	return &LockManager{store: store}
}

func (l *LockManager) Acquire(resource, token string, ttl time.Duration) bool {
	return l.store.SetNX("lock:"+resource, token, ttl)
}

func (l *LockManager) Release(resource, token string) bool {
	return l.store.CompareAndDelete("lock:"+resource, token)
}
