package configcenter

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid gateway config")
	ErrClosed        = errors.New("configuration store closed")
)

type GatewayConfig struct {
	RouteEnabled   bool
	RequestTimeout time.Duration
	RateLimit      int
	RateWindow     time.Duration
	RolloutPercent uint32
	BearerToken    string
}

func (c GatewayConfig) Validate() error {
	switch {
	case c.RequestTimeout <= 0:
		return fmt.Errorf("%w: request timeout must be positive", ErrInvalidConfig)
	case c.RateLimit <= 0:
		return fmt.Errorf("%w: rate limit must be positive", ErrInvalidConfig)
	case c.RateWindow <= 0:
		return fmt.Errorf("%w: rate window must be positive", ErrInvalidConfig)
	case c.RolloutPercent > 100:
		return fmt.Errorf("%w: rollout percent must be between 0 and 100", ErrInvalidConfig)
	case strings.TrimSpace(c.BearerToken) == "":
		return fmt.Errorf("%w: bearer token is required", ErrInvalidConfig)
	default:
		return nil
	}
}

type Snapshot struct {
	Version uint64
	Config  GatewayConfig
}

type Store interface {
	Current() (Snapshot, error)
	Update(GatewayConfig) (Snapshot, error)
	Watch(context.Context) (<-chan Snapshot, error)
	Close() error
}

func InRollout(key string, percent uint32) bool {
	if percent == 0 {
		return false
	}
	if percent >= 100 {
		return true
	}
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return hash.Sum32()%100 < percent
}

func normalizeConfig(config GatewayConfig) GatewayConfig {
	config.BearerToken = strings.TrimSpace(config.BearerToken)
	return config
}
