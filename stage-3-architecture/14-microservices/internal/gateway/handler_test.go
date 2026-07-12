package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
	"just-go/stage-3-architecture/14-microservices/internal/configcenter"
	"just-go/stage-3-architecture/14-microservices/internal/discovery"
)

func TestHandlerReturnsAggregatedProductDetails(t *testing.T) {
	handler, calls := newTestHandler(validGatewayConfig(), nil, nil, nil)
	response := performProductRequest(handler, "teaching-token", "learner-1")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var got productDetails
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	want := productDetails{SKU: "book-1", Name: "Go Microservices", PriceCents: 9900, Quantity: 10, StockVersion: 1}
	if got != want {
		t.Fatalf("response = %#v, want %#v", got, want)
	}
	if calls.product.Load() != 1 || calls.inventory.Load() != 1 {
		t.Fatalf("downstream calls = product:%d inventory:%d", calls.product.Load(), calls.inventory.Load())
	}
}

func TestHandlerRejectsAuthenticationAndRateLimitBeforeDownstream(t *testing.T) {
	config := validGatewayConfig()
	config.RateLimit = 1
	handler, calls := newTestHandler(config, nil, nil, nil)

	unauthorized := performProductRequest(handler, "wrong", "client-a")
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized status = %d", unauthorized.Code)
	}
	first := performProductRequest(handler, "teaching-token", "client-a")
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, body = %s", first.Code, first.Body.String())
	}
	limited := performProductRequest(handler, "teaching-token", "client-a")
	if limited.Code != http.StatusTooManyRequests {
		t.Fatalf("limited status = %d", limited.Code)
	}
	if calls.product.Load() != 1 || calls.inventory.Load() != 1 {
		t.Fatalf("rejected requests called downstream: product:%d inventory:%d", calls.product.Load(), calls.inventory.Load())
	}
}

func TestHandlerAppliesDynamicRouteAndRolloutBeforeDiscovery(t *testing.T) {
	tests := []struct {
		name   string
		config configcenter.GatewayConfig
		key    string
	}{
		{name: "disabled route", config: func() configcenter.GatewayConfig {
			config := validGatewayConfig()
			config.RouteEnabled = false
			return config
		}(), key: "learner-1"},
		{name: "outside rollout", config: func() configcenter.GatewayConfig {
			config := validGatewayConfig()
			config.RolloutPercent = 1
			return config
		}(), key: keyOutsideRollout(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, calls := newTestHandler(tt.config, nil, nil, nil)
			response := performProductRequest(handler, "teaching-token", tt.key)
			if response.Code != http.StatusNotFound {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			if calls.resolver.Load() != 0 || calls.product.Load() != 0 || calls.inventory.Load() != 0 {
				t.Fatalf("disabled route used downstream: %#v", calls)
			}
		})
	}
}

func TestHandlerMapsDiscoveryAndConfigurationFailure(t *testing.T) {
	configFailure := errors.New("config offline")
	handler := NewHandler(
		fakeConfigReader{err: configFailure},
		&fakeResolver{},
		&fakeClientProvider{},
		NewLimiter(nil),
	)
	response := performProductRequest(handler, "teaching-token", "learner-1")
	if response.Code != http.StatusServiceUnavailable || strings.Contains(response.Body.String(), "config offline") {
		t.Fatalf("config failure response = %d %q", response.Code, response.Body.String())
	}

	handler, _ = newTestHandler(validGatewayConfig(), discovery.ErrUnavailable, nil, nil)
	response = performProductRequest(handler, "teaching-token", "learner-2")
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("discovery failure status = %d", response.Code)
	}
}

func TestHandlerMapsDownstreamErrorsWithoutLeakingDetails(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "invalid", err: status.Error(codes.InvalidArgument, "secret invalid detail"), want: http.StatusBadRequest},
		{name: "missing", err: status.Error(codes.NotFound, "secret missing detail"), want: http.StatusNotFound},
		{name: "unavailable", err: status.Error(codes.Unavailable, "secret unavailable detail"), want: http.StatusServiceUnavailable},
		{name: "deadline", err: status.Error(codes.DeadlineExceeded, "secret timeout detail"), want: http.StatusGatewayTimeout},
		{name: "unknown", err: errors.New("secret internal detail"), want: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productClient := &fakeProductClient{err: tt.err}
			handler, _ := newTestHandler(validGatewayConfig(), nil, productClient, nil)
			response := performProductRequest(handler, "teaching-token", "key-"+tt.name)
			if response.Code != tt.want {
				t.Fatalf("status = %d, want %d, body = %s", response.Code, tt.want, response.Body.String())
			}
			body := response.Body.String()
			if strings.Contains(body, "secret") || strings.Contains(body, "book-1") {
				t.Fatalf("error response leaked detail or partial data: %q", body)
			}
		})
	}
}

func TestHandlerCallsDownstreamsConcurrently(t *testing.T) {
	productStarted := make(chan struct{})
	inventoryStarted := make(chan struct{})
	release := make(chan struct{})
	productClient := &fakeProductClient{call: func(ctx context.Context, _ *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
		close(productStarted)
		select {
		case <-release:
			return defaultProductResponse(), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}}
	inventoryClient := &fakeInventoryClient{call: func(ctx context.Context, _ *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error) {
		close(inventoryStarted)
		select {
		case <-release:
			return defaultStockResponse(), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}}
	handler, _ := newTestHandler(validGatewayConfig(), nil, productClient, inventoryClient)
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		done <- performProductRequest(handler, "teaching-token", "concurrent-key")
	}()

	select {
	case <-productStarted:
	case <-time.After(time.Second):
		t.Fatal("product call did not start")
	}
	select {
	case <-inventoryStarted:
	case <-time.After(time.Second):
		t.Fatal("inventory call did not start before product was released")
	}
	close(release)
	if response := <-done; response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestHandlerEnforcesConfiguredDeadline(t *testing.T) {
	config := validGatewayConfig()
	config.RequestTimeout = 20 * time.Millisecond
	productClient := &fakeProductClient{call: func(ctx context.Context, _ *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	handler, _ := newTestHandler(config, nil, productClient, nil)
	response := performProductRequest(handler, "teaching-token", "deadline-key")
	if response.Code != http.StatusGatewayTimeout {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

type callCounts struct {
	resolver  atomic.Int32
	product   atomic.Int32
	inventory atomic.Int32
}

type fakeConfigReader struct {
	snapshot configcenter.Snapshot
	err      error
}

func (f fakeConfigReader) Current() (configcenter.Snapshot, error) { return f.snapshot, f.err }

type fakeResolver struct {
	calls *atomic.Int32
	err   error
}

func (f *fakeResolver) Resolve(service string) (discovery.Instance, error) {
	if f.calls != nil {
		f.calls.Add(1)
	}
	if f.err != nil {
		return discovery.Instance{}, f.err
	}
	port := "5001"
	if service == "inventory" {
		port = "5002"
	}
	return discovery.Instance{Service: service, ID: service + "-1", Address: "127.0.0.1:" + port}, nil
}

type fakeClientProvider struct {
	calls     *callCounts
	product   productv1.ProductServiceClient
	inventory inventoryv1.InventoryServiceClient
}

func (f *fakeClientProvider) Product(context.Context, string) (productv1.ProductServiceClient, error) {
	if f.calls != nil {
		f.calls.product.Add(1)
	}
	return f.product, nil
}

func (f *fakeClientProvider) Inventory(context.Context, string) (inventoryv1.InventoryServiceClient, error) {
	if f.calls != nil {
		f.calls.inventory.Add(1)
	}
	return f.inventory, nil
}

type fakeProductClient struct {
	call func(context.Context, *productv1.GetProductRequest) (*productv1.GetProductResponse, error)
	err  error
}

func (f *fakeProductClient) GetProduct(ctx context.Context, req *productv1.GetProductRequest, _ ...grpc.CallOption) (*productv1.GetProductResponse, error) {
	if f.call != nil {
		return f.call(ctx, req)
	}
	if f.err != nil {
		return nil, f.err
	}
	return defaultProductResponse(), nil
}

type fakeInventoryClient struct {
	call func(context.Context, *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error)
	err  error
}

func (f *fakeInventoryClient) GetStock(ctx context.Context, req *inventoryv1.GetStockRequest, _ ...grpc.CallOption) (*inventoryv1.GetStockResponse, error) {
	if f.call != nil {
		return f.call(ctx, req)
	}
	if f.err != nil {
		return nil, f.err
	}
	return defaultStockResponse(), nil
}

func (f *fakeInventoryClient) WatchStock(context.Context, *inventoryv1.WatchStockRequest, ...grpc.CallOption) (grpc.ServerStreamingClient[inventoryv1.WatchStockResponse], error) {
	return nil, status.Error(codes.Unimplemented, "unused in gateway")
}

func (f *fakeInventoryClient) SyncStock(context.Context, ...grpc.CallOption) (grpc.BidiStreamingClient[inventoryv1.SyncStockRequest, inventoryv1.SyncStockResponse], error) {
	return nil, status.Error(codes.Unimplemented, "unused in gateway")
}

func newTestHandler(config configcenter.GatewayConfig, resolveErr error, productClient *fakeProductClient, inventoryClient *fakeInventoryClient) (http.Handler, *callCounts) {
	if productClient == nil {
		productClient = &fakeProductClient{}
	}
	if inventoryClient == nil {
		inventoryClient = &fakeInventoryClient{}
	}
	calls := &callCounts{}
	return NewHandler(
		fakeConfigReader{snapshot: configcenter.Snapshot{Version: 1, Config: config}},
		&fakeResolver{calls: &calls.resolver, err: resolveErr},
		&fakeClientProvider{calls: calls, product: productClient, inventory: inventoryClient},
		NewLimiter(func() time.Time { return time.Unix(100, 0) }),
	), calls
}

func performProductRequest(handler http.Handler, token, requestKey string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/products/book-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Request-Key", requestKey)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func validGatewayConfig() configcenter.GatewayConfig {
	return configcenter.GatewayConfig{
		RouteEnabled:   true,
		RequestTimeout: time.Second,
		RateLimit:      100,
		RateWindow:     time.Minute,
		RolloutPercent: 100,
		BearerToken:    "teaching-token",
	}
}

func defaultProductResponse() *productv1.GetProductResponse {
	return &productv1.GetProductResponse{Sku: "book-1", Name: "Go Microservices", PriceCents: 9900}
}

func defaultStockResponse() *inventoryv1.GetStockResponse {
	return &inventoryv1.GetStockResponse{Sku: "book-1", Quantity: 10, Version: 1}
}

func keyOutsideRollout(percent uint32) string {
	for i := range 10_000 {
		key := "outside-" + time.Unix(int64(i), 0).Format(time.RFC3339Nano)
		if !configcenter.InRollout(key, percent) {
			return key
		}
	}
	panic("could not find rollout key")
}
