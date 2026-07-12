package discovery

import (
	"context"
	"errors"
)

var (
	ErrInvalidInstance   = errors.New("invalid service instance")
	ErrDuplicateInstance = errors.New("duplicate service instance")
	ErrUnavailable       = errors.New("service unavailable")
	ErrClosed            = errors.New("service discovery closed")
)

type Instance struct {
	Service string
	ID      string
	Address string
}

type Registry interface {
	Register(Instance) (deregister func() error, err error)
	Resolve(service string) (Instance, error)
	Watch(ctx context.Context, service string) (<-chan []Instance, error)
	Close() error
}
