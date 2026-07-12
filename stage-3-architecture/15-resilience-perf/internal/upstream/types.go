package upstream

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type Product struct {
	SKU           string `json:"sku"`
	Name          string `json:"name"`
	PriceCents    int64  `json:"price_cents"`
	Quantity      int64  `json:"quantity"`
	Degraded      bool   `json:"degraded"`
	DegradeReason string `json:"degrade_reason,omitempty"`
}

type transportError struct {
	err error
}

func (e transportError) Error() string { return e.err.Error() }
func (e transportError) Unwrap() error { return e.err }

type Error struct {
	StatusCode int
	Temporary  bool
	Message    string
}

func (e Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("upstream status %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("upstream status %d", e.StatusCode)
}

func IsRetryable(err error) bool {
	var upstreamErr Error
	if errors.As(err, &upstreamErr) {
		return upstreamErr.Temporary || upstreamErr.StatusCode == http.StatusTooManyRequests || upstreamErr.StatusCode >= 500
	}
	var transportErr transportError
	if errors.As(err, &transportErr) {
		return true
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}

func IsClientError(err error) bool {
	var upstreamErr Error
	if !errors.As(err, &upstreamErr) {
		return false
	}
	return upstreamErr.StatusCode >= 400 && upstreamErr.StatusCode < 500 && upstreamErr.StatusCode != http.StatusTooManyRequests
}

func IsRateLimited(err error) bool {
	var upstreamErr Error
	return errors.As(err, &upstreamErr) && upstreamErr.StatusCode == http.StatusTooManyRequests
}

func IsNotFound(err error) bool {
	var upstreamErr Error
	return errors.As(err, &upstreamErr) && upstreamErr.StatusCode == http.StatusNotFound
}
