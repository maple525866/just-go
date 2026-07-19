package breaker

import (
	"errors"
	"time"

	"github.com/sony/gobreaker/v2"
)

var ErrOpen = errors.New("circuit breaker open")

type Config struct {
	Name             string
	MaxRequests      uint32
	Timeout          time.Duration
	FailureThreshold uint32
	IsExcluded       func(error) bool
}

type Circuit[T any] struct {
	cb *gobreaker.CircuitBreaker[T]
}

func New[T any](config Config) *Circuit[T] {
	if config.Name == "" {
		config.Name = "chapter-15"
	}
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}
	if config.Timeout <= 0 {
		config.Timeout = time.Second
	}
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 3
	}
	settings := gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Timeout:     config.Timeout,
		IsExcluded:  config.IsExcluded,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= config.FailureThreshold
		},
	}
	return &Circuit[T]{cb: gobreaker.NewCircuitBreaker[T](settings)}
}

func (c *Circuit[T]) Execute(fn func() (T, error)) (T, error) {
	result, err := c.cb.Execute(fn)
	if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
		return result, ErrOpen
	}
	return result, err
}

func (c *Circuit[T]) State() string {
	switch c.cb.State() {
	case gobreaker.StateClosed:
		return "closed"
	case gobreaker.StateHalfOpen:
		return "half-open"
	case gobreaker.StateOpen:
		return "open"
	default:
		return "unknown"
	}
}
