package gateway

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
	"just-go/stage-3-architecture/14-microservices/internal/configcenter"
	"just-go/stage-3-architecture/14-microservices/internal/discovery"
	inventoryservice "just-go/stage-3-architecture/14-microservices/internal/inventory"
	productservice "just-go/stage-3-architecture/14-microservices/internal/product"
)

func TestGatewayIntegration(t *testing.T) {
	catalog, err := productservice.NewCatalog([]productservice.Product{{
		SKU: "book-1", Name: "Go Microservices", PriceCents: 9900,
	}})
	if err != nil {
		t.Fatal(err)
	}
	stock, err := inventoryservice.NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	productAddress := startGRPCServer(t, func(server *grpc.Server) {
		productv1.RegisterProductServiceServer(server, productservice.NewService(catalog))
	})
	inventoryAddress := startGRPCServer(t, func(server *grpc.Server) {
		inventoryv1.RegisterInventoryServiceServer(server, inventoryservice.NewService(stock))
	})

	registry := discovery.NewMemoryRegistry()
	t.Cleanup(func() { _ = registry.Close() })
	if _, err := registry.Register(discovery.Instance{Service: "product", ID: "product-1", Address: productAddress}); err != nil {
		t.Fatal(err)
	}
	deregisterInventory, err := registry.Register(discovery.Instance{Service: "inventory", ID: "inventory-1", Address: inventoryAddress})
	if err != nil {
		t.Fatal(err)
	}
	configStore, err := configcenter.NewMemoryStore(validGatewayConfig())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = configStore.Close() })
	connections := NewConnections(DefaultDial)
	t.Cleanup(func() { _ = connections.Close() })
	httpServer := httptest.NewServer(NewHandler(configStore, registry, connections, NewLimiter(nil)))
	t.Cleanup(httpServer.Close)

	statusCode, body := integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "success-key")
	if statusCode != http.StatusOK || !strings.Contains(body, `"name":"Go Microservices"`) || !strings.Contains(body, `"quantity":10`) {
		t.Fatalf("success = %d %s", statusCode, body)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "wrong", "auth-key")
	if statusCode != http.StatusUnauthorized {
		t.Fatalf("unauthorized status = %d", statusCode)
	}

	config := validGatewayConfig()
	config.RouteEnabled = false
	if _, err := configStore.Update(config); err != nil {
		t.Fatal(err)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "disabled-key")
	if statusCode != http.StatusNotFound {
		t.Fatalf("disabled status = %d", statusCode)
	}
	config.RouteEnabled = true
	config.RolloutPercent = 0
	if _, err := configStore.Update(config); err != nil {
		t.Fatal(err)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "rollout-key")
	if statusCode != http.StatusNotFound {
		t.Fatalf("rollout status = %d", statusCode)
	}
	config.RolloutPercent = 100
	config.RateLimit = 1
	if _, err := configStore.Update(config); err != nil {
		t.Fatal(err)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "limited-key")
	if statusCode != http.StatusOK {
		t.Fatalf("first limited-key status = %d", statusCode)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "limited-key")
	if statusCode != http.StatusTooManyRequests {
		t.Fatalf("second limited-key status = %d", statusCode)
	}

	config.RateLimit = 100
	if _, err := configStore.Update(config); err != nil {
		t.Fatal(err)
	}
	if err := deregisterInventory(); err != nil {
		t.Fatal(err)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "discovery-key")
	if statusCode != http.StatusServiceUnavailable {
		t.Fatalf("discovery status = %d", statusCode)
	}
	deregisterInventory, err = registry.Register(discovery.Instance{Service: "inventory", ID: "inventory-2", Address: inventoryAddress})
	if err != nil {
		t.Fatal(err)
	}
	statusCode, body = integrationRequest(t, httpServer.Client(), httpServer.URL, "missing", "teaching-token", "missing-key")
	if statusCode != http.StatusNotFound || strings.Contains(body, "secret") {
		t.Fatalf("missing response = %d %s", statusCode, body)
	}

	if err := deregisterInventory(); err != nil {
		t.Fatal(err)
	}
	slowAddress := startGRPCServer(t, func(server *grpc.Server) {
		inventoryv1.RegisterInventoryServiceServer(server, slowInventoryService{})
	})
	if _, err := registry.Register(discovery.Instance{Service: "inventory", ID: "inventory-slow", Address: slowAddress}); err != nil {
		t.Fatal(err)
	}
	config.RequestTimeout = 20 * time.Millisecond
	if _, err := configStore.Update(config); err != nil {
		t.Fatal(err)
	}
	statusCode, _ = integrationRequest(t, httpServer.Client(), httpServer.URL, "book-1", "teaching-token", "timeout-key")
	if statusCode != http.StatusGatewayTimeout {
		t.Fatalf("timeout status = %d", statusCode)
	}
}

type slowInventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
}

func (slowInventoryService) GetStock(ctx context.Context, _ *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error) {
	<-ctx.Done()
	return nil, status.FromContextError(ctx.Err()).Err()
}

func startGRPCServer(t *testing.T, register func(*grpc.Server)) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	register(server)
	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})
	return listener.Addr().String()
}

func integrationRequest(t *testing.T, client *http.Client, baseURL, sku, token, key string) (int, string) {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, baseURL+"/api/v1/products/"+sku, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Request-Key", key)
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, string(body)
}
