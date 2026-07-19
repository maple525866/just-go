package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"just-go/stage-3-architecture/15-resilience-perf/internal/breaker"
	"just-go/stage-3-architecture/15-resilience-perf/internal/bulkhead"
	"just-go/stage-3-architecture/15-resilience-perf/internal/limiter"
	"just-go/stage-3-architecture/15-resilience-perf/internal/retry"
	"just-go/stage-3-architecture/15-resilience-perf/internal/upstream"
)

type staticLimiter struct{ decision limiter.Decision }

func (l staticLimiter) Allow() limiter.Decision { return l.decision }

type fakeClient struct {
	calls    int
	products []upstream.Product
	errors   []error
}

func (c *fakeClient) GetProduct(ctx context.Context, sku string) (upstream.Product, error) {
	c.calls++
	index := c.calls - 1
	if index < len(c.errors) && c.errors[index] != nil {
		return upstream.Product{}, c.errors[index]
	}
	if index < len(c.products) {
		return c.products[index], nil
	}
	return upstream.Product{SKU: sku, Name: "Go Resilience", PriceCents: 9900, Quantity: 5}, nil
}

type passCircuit struct{}

func (passCircuit) Execute(fn func() (upstream.Product, error)) (upstream.Product, error) {
	return fn()
}
func (passCircuit) State() string { return "closed" }

type openCircuit struct{}

func (openCircuit) Execute(fn func() (upstream.Product, error)) (upstream.Product, error) {
	return upstream.Product{}, breaker.ErrOpen
}
func (openCircuit) State() string { return "open" }

func newTestBulkhead(t *testing.T, limit int) *bulkhead.Bulkhead {
	t.Helper()
	b, err := bulkhead.New(limit)
	if err != nil {
		t.Fatalf("bulkhead.New: %v", err)
	}
	return b
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v body=%s", err, recorder.Body.String())
	}
	return body
}

func TestHandlerSuccess(t *testing.T) {
	client := &fakeClient{products: []upstream.Product{{SKU: "book-1", Name: "Go Resilience", PriceCents: 9900, Quantity: 8}}}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true, Remaining: 9}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 1},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	body := decodeResponse(t, rr)
	if body["sku"] != "book-1" || body["degraded"] != false || body["attempts"].(float64) != 1 {
		t.Fatalf("body=%#v", body)
	}
	if client.calls != 1 {
		t.Fatalf("calls=%d", client.calls)
	}
}

func TestHandlerRateLimitRejectsBeforeUpstream(t *testing.T) {
	client := &fakeClient{}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: false, RetryAfter: 250 * time.Millisecond}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 1},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	body := decodeResponse(t, rr)
	if body["retry_after_millis"].(float64) != 250 {
		t.Fatalf("body=%#v", body)
	}
	if client.calls != 0 {
		t.Fatalf("upstream called %d times", client.calls)
	}
}

func TestHandlerBulkheadFull(t *testing.T) {
	b := newTestBulkhead(t, 1)
	release, err := b.Acquire(context.Background())
	if err != nil {
		t.Fatalf("pre-acquire: %v", err)
	}
	defer release()

	handler, err := NewHandler(Options{
		Client:          &fakeClient{},
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        b,
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 1},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	body := decodeResponse(t, rr)
	if body["reason"] != "bulkhead_full" {
		t.Fatalf("body=%#v", body)
	}
}

func TestHandlerRetriesThenSucceeds(t *testing.T) {
	client := &fakeClient{
		errors:   []error{upstream.Error{StatusCode: http.StatusServiceUnavailable, Temporary: true, Message: "try again"}, nil},
		products: []upstream.Product{{}, {SKU: "book-1", Name: "Recovered", PriceCents: 1, Quantity: 2}},
	}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 2, BaseDelay: time.Millisecond, Sleep: func(context.Context, time.Duration) error { return nil }, ShouldRetry: upstream.IsRetryable},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	body := decodeResponse(t, rr)
	if body["name"] != "Recovered" || body["attempts"].(float64) != 2 {
		t.Fatalf("body=%#v", body)
	}
}

func TestHandlerCircuitOpenUsesFallback(t *testing.T) {
	handler, err := NewHandler(Options{
		Client:          &fakeClient{},
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         openCircuit{},
		Retry:           retry.Policy{MaxAttempts: 1},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	body := decodeResponse(t, rr)
	if body["degraded"] != true || body["degrade_reason"] != "circuit_open" {
		t.Fatalf("body=%#v", body)
	}
}

func TestHandlerClientErrorDoesNotFallback(t *testing.T) {
	client := &fakeClient{errors: []error{upstream.Error{StatusCode: http.StatusNotFound, Message: "missing"}}}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 3, ShouldRetry: upstream.IsRetryable},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/missing", http.NoBody))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if client.calls != 1 {
		t.Fatalf("client errors must not retry, calls=%d", client.calls)
	}
}

func TestHandlerTimeoutReturnsGatewayTimeoutWhenFallbackDisabled(t *testing.T) {
	client := &fakeClient{errors: []error{context.DeadlineExceeded}}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 1},
		Timeout:         time.Nanosecond,
		FallbackEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
	if rr.Code != http.StatusGatewayTimeout {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestHandlerUpstreamRateLimitPreservesBackpressure(t *testing.T) {
	client := &fakeClient{errors: []error{
		upstream.Error{StatusCode: http.StatusTooManyRequests, Temporary: true, Message: "slow down"},
		upstream.Error{StatusCode: http.StatusTooManyRequests, Temporary: true, Message: "slow down"},
	}}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 2, BaseDelay: time.Millisecond, Sleep: func(context.Context, time.Duration) error { return nil }, ShouldRetry: upstream.IsRetryable},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if client.calls != 2 {
		t.Fatalf("expected two retry attempts, got %d", client.calls)
	}
}

func TestHandlerRejectsNestedProductPathsBeforePolicies(t *testing.T) {
	client := &fakeClient{}
	handler, err := NewHandler(Options{
		Client:          client,
		Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
		Bulkhead:        newTestBulkhead(t, 1),
		Circuit:         passCircuit{},
		Retry:           retry.Policy{MaxAttempts: 1},
		Timeout:         time.Second,
		FallbackEnabled: true,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1/reviews", http.NoBody))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if client.calls != 0 {
		t.Fatalf("nested route should not call upstream, calls=%d", client.calls)
	}
}

func TestHandlerBreakerExcludesClientErrorsAndCanceledRequests(t *testing.T) {
	tests := []struct {
		name       string
		excluded   error
		wantStatus int
	}{
		{
			name:       "upstream client error",
			excluded:   upstream.Error{StatusCode: http.StatusBadRequest, Message: "bad sku"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "upstream rate limit",
			excluded:   upstream.Error{StatusCode: http.StatusTooManyRequests, Temporary: true, Message: "slow down"},
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name:       "request canceled",
			excluded:   context.Canceled,
			wantStatus: http.StatusRequestTimeout,
		},
		{
			name:       "deadline exceeded",
			excluded:   context.DeadlineExceeded,
			wantStatus: http.StatusGatewayTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fakeClient{
				errors: []error{tt.excluded, tt.excluded, nil},
				products: []upstream.Product{
					{},
					{},
					{SKU: "book-1", Name: "Recovered", PriceCents: 1200, Quantity: 7},
				},
			}
			handler, err := NewHandler(Options{
				Client:          client,
				Limiter:         staticLimiter{decision: limiter.Decision{Allowed: true}},
				Bulkhead:        newTestBulkhead(t, 1),
				Circuit:         breaker.New[upstream.Product](breaker.Config{Name: tt.name, FailureThreshold: 2, Timeout: time.Second, MaxRequests: 1, IsExcluded: IsBreakerExcluded}),
				Retry:           retry.Policy{MaxAttempts: 1},
				Timeout:         time.Second,
				FallbackEnabled: true,
			})
			if err != nil {
				t.Fatalf("NewHandler: %v", err)
			}

			for i := 0; i < 2; i++ {
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
				if rr.Code != tt.wantStatus {
					t.Fatalf("excluded request %d status=%d body=%s", i+1, rr.Code, rr.Body.String())
				}
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", http.NoBody))
			if rr.Code != http.StatusOK {
				t.Fatalf("valid request status=%d body=%s", rr.Code, rr.Body.String())
			}
			body := decodeResponse(t, rr)
			if body["degraded"] != false || body["name"] != "Recovered" {
				t.Fatalf("valid request should use upstream, body=%#v", body)
			}
		})
	}
}
