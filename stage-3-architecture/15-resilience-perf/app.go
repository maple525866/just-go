package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"just-go/stage-3-architecture/15-resilience-perf/internal/breaker"
	"just-go/stage-3-architecture/15-resilience-perf/internal/bulkhead"
	"just-go/stage-3-architecture/15-resilience-perf/internal/gateway"
	"just-go/stage-3-architecture/15-resilience-perf/internal/limiter"
	"just-go/stage-3-architecture/15-resilience-perf/internal/profiler"
	"just-go/stage-3-architecture/15-resilience-perf/internal/retry"
	"just-go/stage-3-architecture/15-resilience-perf/internal/upstream"
)

func RunDemo(ctx context.Context, out io.Writer) error {
	upstreamHandler := upstream.NewScriptedHandler([]upstream.ScriptedResponse{{
		StatusCode: http.StatusOK,
		Product:    upstream.Product{SKU: "book-1", Name: "Go Resilience", PriceCents: 9900, Quantity: 10},
	}, {
		StatusCode: http.StatusServiceUnavailable,
		Body:       `{"error":"temporary"}`,
	}, {
		StatusCode: http.StatusServiceUnavailable,
		Body:       `{"error":"still temporary"}`,
	}})
	upstreamServer := httptest.NewServer(upstreamHandler)
	defer upstreamServer.Close()

	bucket, err := limiter.NewTokenBucket(limiter.Config{Capacity: 2, RefillPerSecond: 0.000001}, nil)
	if err != nil {
		return err
	}
	bh, err := bulkhead.New(2)
	if err != nil {
		return err
	}
	circuit := breaker.New[upstream.Product](breaker.Config{Name: "demo-product", FailureThreshold: 2, Timeout: 100 * time.Millisecond, MaxRequests: 1, IsExcluded: gateway.IsBreakerExcluded})
	handler, err := gateway.NewHandler(gateway.Options{
		Client:          upstream.NewHTTPClient(upstreamServer.URL, upstreamServer.Client()),
		Limiter:         bucket,
		Bulkhead:        bh,
		Circuit:         circuit,
		Retry:           retry.Policy{MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, Sleep: func(context.Context, time.Duration) error { return nil }, ShouldRetry: upstream.IsRetryable},
		Timeout:         500 * time.Millisecond,
		FallbackEnabled: true,
	})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/products/", handler)
	profiler.Register(mux)
	gatewayServer := httptest.NewServer(mux)
	defer gatewayServer.Close()

	if err := requestDemo(ctx, out, gatewayServer.Client(), gatewayServer.URL+"/api/v1/products/book-1", "success"); err != nil {
		return err
	}
	if err := requestDemo(ctx, out, gatewayServer.Client(), gatewayServer.URL+"/api/v1/products/book-1", "fallback"); err != nil {
		return err
	}
	if err := requestDemo(ctx, out, gatewayServer.Client(), gatewayServer.URL+"/api/v1/products/book-1", "rate-limit"); err != nil {
		return err
	}
	fmt.Fprintf(out, "heap-demo bytes=%d\n", profiler.AllocateHotHeap(2, 1024))
	return nil
}

func requestDemo(ctx context.Context, out io.Writer, client *http.Client, url string, label string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Fprintf(out, "%s status=%d\n", label, resp.StatusCode)
	return nil
}
