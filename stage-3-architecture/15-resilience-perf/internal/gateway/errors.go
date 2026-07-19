package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"just-go/stage-3-architecture/15-resilience-perf/internal/breaker"
	"just-go/stage-3-architecture/15-resilience-perf/internal/bulkhead"
	"just-go/stage-3-architecture/15-resilience-perf/internal/upstream"
)

type errorResponse struct {
	Error            string `json:"error"`
	Reason           string `json:"reason,omitempty"`
	RetryAfterMillis int64  `json:"retry_after_millis,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func IsBreakerExcluded(err error) bool {
	return upstream.IsClientError(err) || upstream.IsRateLimited(err) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func statusForError(err error) (int, string) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout, "deadline_exceeded"
	case errors.Is(err, context.Canceled):
		return http.StatusRequestTimeout, "request_canceled"
	case errors.Is(err, bulkhead.ErrFull):
		return http.StatusServiceUnavailable, "bulkhead_full"
	case errors.Is(err, breaker.ErrOpen):
		return http.StatusServiceUnavailable, "circuit_open"
	case upstream.IsNotFound(err):
		return http.StatusNotFound, "not_found"
	case upstream.IsRateLimited(err):
		return http.StatusTooManyRequests, "upstream_rate_limited"
	case upstream.IsClientError(err):
		return http.StatusBadRequest, "upstream_client_error"
	case upstream.IsRetryable(err):
		return http.StatusServiceUnavailable, "upstream_unavailable"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
