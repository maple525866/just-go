package gateway

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"just-go/stage-3-architecture/15-resilience-perf/internal/breaker"
	"just-go/stage-3-architecture/15-resilience-perf/internal/bulkhead"
	"just-go/stage-3-architecture/15-resilience-perf/internal/limiter"
	"just-go/stage-3-architecture/15-resilience-perf/internal/retry"
	"just-go/stage-3-architecture/15-resilience-perf/internal/upstream"
)

type RateLimiter interface {
	Allow() limiter.Decision
}

type Circuit interface {
	Execute(func() (upstream.Product, error)) (upstream.Product, error)
	State() string
}

type Options struct {
	Client          upstream.Client
	Limiter         RateLimiter
	Bulkhead        *bulkhead.Bulkhead
	Circuit         Circuit
	Retry           retry.Policy
	Timeout         time.Duration
	FallbackEnabled bool
}

type handler struct{ options Options }

func NewHandler(options Options) (http.Handler, error) {
	if options.Client == nil {
		return nil, errors.New("client is required")
	}
	if options.Limiter == nil {
		return nil, errors.New("limiter is required")
	}
	if options.Bulkhead == nil {
		return nil, errors.New("bulkhead is required")
	}
	if options.Circuit == nil {
		return nil, errors.New("circuit is required")
	}
	if options.Timeout <= 0 {
		options.Timeout = time.Second
	}
	return &handler{options: options}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !strings.HasPrefix(r.URL.Path, "/api/v1/products/") {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not_found", Reason: "route_not_found"})
		return
	}
	sku := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	if sku == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "bad_request", Reason: "sku_required"})
		return
	}
	if strings.Contains(sku, "/") {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not_found", Reason: "route_not_found"})
		return
	}

	decision := h.options.Limiter.Allow()
	if !decision.Allowed {
		writeJSON(w, http.StatusTooManyRequests, errorResponse{Error: "rate_limited", Reason: "token_bucket_empty", RetryAfterMillis: decision.RetryAfter.Milliseconds()})
		return
	}

	release, err := h.options.Bulkhead.Acquire(r.Context())
	if err != nil {
		status, reason := statusForError(err)
		writeJSON(w, status, errorResponse{Error: "unavailable", Reason: reason})
		return
	}
	defer release()

	ctx, cancel := context.WithTimeout(r.Context(), h.options.Timeout)
	defer cancel()

	product, stats, err := retry.Do[upstream.Product](ctx, h.options.Retry, func(ctx context.Context) (upstream.Product, error) {
		return h.options.Circuit.Execute(func() (upstream.Product, error) {
			return h.options.Client.GetProduct(ctx, sku)
		})
	})
	if err == nil {
		writeJSON(w, http.StatusOK, productResponse(product, stats.Attempts))
		return
	}

	if h.options.FallbackEnabled && !upstream.IsClientError(err) && !upstream.IsRateLimited(err) && !upstream.IsNotFound(err) && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		fallback := upstream.Product{SKU: sku, Name: "temporarily unavailable", PriceCents: 0, Quantity: 0, Degraded: true, DegradeReason: fallbackReason(err)}
		writeJSON(w, http.StatusOK, productResponse(fallback, stats.Attempts))
		return
	}

	status, reason := statusForError(err)
	writeJSON(w, status, errorResponse{Error: http.StatusText(status), Reason: reason})
}

type response struct {
	SKU           string `json:"sku"`
	Name          string `json:"name"`
	PriceCents    int64  `json:"price_cents"`
	Quantity      int64  `json:"quantity"`
	Degraded      bool   `json:"degraded"`
	DegradeReason string `json:"degrade_reason,omitempty"`
	Attempts      int    `json:"attempts"`
}

func productResponse(product upstream.Product, attempts int) response {
	return response{SKU: product.SKU, Name: product.Name, PriceCents: product.PriceCents, Quantity: product.Quantity, Degraded: product.Degraded, DegradeReason: product.DegradeReason, Attempts: attempts}
}

func fallbackReason(err error) string {
	switch {
	case errors.Is(err, breaker.ErrOpen):
		return "circuit_open"
	case upstream.IsRetryable(err):
		return "upstream_unavailable"
	default:
		return "fallback"
	}
}
