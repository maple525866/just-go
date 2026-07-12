package configcenter

import (
	"errors"
	"testing"
	"time"
)

func TestGatewayConfigValidate(t *testing.T) {
	valid := validConfig()
	tests := []struct {
		name   string
		mutate func(*GatewayConfig)
	}{
		{name: "zero timeout", mutate: func(c *GatewayConfig) { c.RequestTimeout = 0 }},
		{name: "negative timeout", mutate: func(c *GatewayConfig) { c.RequestTimeout = -time.Second }},
		{name: "zero rate limit", mutate: func(c *GatewayConfig) { c.RateLimit = 0 }},
		{name: "zero rate window", mutate: func(c *GatewayConfig) { c.RateWindow = 0 }},
		{name: "rollout over 100", mutate: func(c *GatewayConfig) { c.RolloutPercent = 101 }},
		{name: "blank token", mutate: func(c *GatewayConfig) { c.BearerToken = " " }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := valid
			tt.mutate(&config)
			if err := config.Validate(); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("error = %v, want ErrInvalidConfig", err)
			}
		})
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid config: %v", err)
	}
}

func TestInRolloutIsStableAndBounded(t *testing.T) {
	first := InRollout("learner-42", 25)
	for range 100 {
		if got := InRollout("learner-42", 25); got != first {
			t.Fatalf("decision changed: first=%v got=%v", first, got)
		}
	}
	if InRollout("any", 0) {
		t.Fatal("zero-percent rollout enabled")
	}
	if !InRollout("any", 100) {
		t.Fatal("full rollout disabled")
	}
}

func validConfig() GatewayConfig {
	return GatewayConfig{
		RouteEnabled:   true,
		RequestTimeout: time.Second,
		RateLimit:      10,
		RateWindow:     time.Minute,
		RolloutPercent: 100,
		BearerToken:    "teaching-token",
	}
}
