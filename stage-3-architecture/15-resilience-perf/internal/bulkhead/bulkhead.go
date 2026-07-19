package bulkhead

import (
	"context"
	"errors"
	"sync"
)

var ErrFull = errors.New("bulkhead full")

type Bulkhead struct {
	slots chan struct{}
}

func New(limit int) (*Bulkhead, error) {
	if limit <= 0 {
		return nil, errors.New("bulkhead limit must be positive")
	}
	return &Bulkhead{slots: make(chan struct{}, limit)}, nil
}

func (b *Bulkhead) Acquire(ctx context.Context) (func(), error) {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	select {
	case b.slots <- struct{}{}:
		var once sync.Once
		return func() {
			once.Do(func() { <-b.slots })
		}, nil
	default:
		return nil, ErrFull
	}
}
