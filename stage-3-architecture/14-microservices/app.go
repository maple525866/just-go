package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
	"just-go/stage-3-architecture/14-microservices/internal/configcenter"
	"just-go/stage-3-architecture/14-microservices/internal/discovery"
	"just-go/stage-3-architecture/14-microservices/internal/gateway"
	inventoryservice "just-go/stage-3-architecture/14-microservices/internal/inventory"
	productservice "just-go/stage-3-architecture/14-microservices/internal/product"
)

func run(ctx context.Context, out io.Writer) (retErr error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	if out == nil {
		return errors.New("output writer is required")
	}

	var (
		registry            *discovery.MemoryRegistry
		configStore         *configcenter.MemoryStore
		connections         *gateway.Connections
		productServer       *grpc.Server
		inventoryServer     *grpc.Server
		httpServer          *http.Server
		deregisterProduct   func() error
		deregisterInventory func() error
		serveWG             sync.WaitGroup
	)
	defer func() {
		var cleanupErrors []error
		if httpServer != nil {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("shutdown gateway: %w", err))
			}
			cancel()
		}
		if connections != nil {
			if err := connections.Close(); err != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("close grpc clients: %w", err))
			}
		}
		if productServer != nil {
			productServer.GracefulStop()
		}
		if inventoryServer != nil {
			inventoryServer.GracefulStop()
		}
		if deregisterInventory != nil {
			if err := deregisterInventory(); err != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("deregister inventory: %w", err))
			}
		}
		if deregisterProduct != nil {
			if err := deregisterProduct(); err != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("deregister product: %w", err))
			}
		}
		if configStore != nil {
			if err := configStore.Close(); err != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("close config store: %w", err))
			}
		}
		if registry != nil {
			if err := registry.Close(); err != nil {
				cleanupErrors = append(cleanupErrors, fmt.Errorf("close discovery: %w", err))
			}
		}
		serveWG.Wait()
		retErr = errors.Join(retErr, errors.Join(cleanupErrors...))
	}()

	catalog, err := productservice.NewCatalog([]productservice.Product{{
		SKU: "book-1", Name: "Go Microservices", PriceCents: 9900,
	}})
	if err != nil {
		return fmt.Errorf("create product catalog: %w", err)
	}
	stock, err := inventoryservice.NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		return fmt.Errorf("create inventory store: %w", err)
	}
	registry = discovery.NewMemoryRegistry()
	configStore, err = configcenter.NewMemoryStore(configcenter.GatewayConfig{
		RouteEnabled:   true,
		RequestTimeout: time.Second,
		RateLimit:      10,
		RateWindow:     time.Minute,
		RolloutPercent: 100,
		BearerToken:    "teaching-token",
	})
	if err != nil {
		return fmt.Errorf("create config store: %w", err)
	}

	productListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen product grpc: %w", err)
	}
	productServer = grpc.NewServer()
	productv1.RegisterProductServiceServer(productServer, productservice.NewService(catalog))
	serveWG.Add(1)
	go func() {
		defer serveWG.Done()
		_ = productServer.Serve(productListener)
	}()
	deregisterProduct, err = registry.Register(discovery.Instance{
		Service: "product", ID: "product-demo", Address: productListener.Addr().String(),
	})
	if err != nil {
		return fmt.Errorf("register product: %w", err)
	}

	inventoryListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen inventory grpc: %w", err)
	}
	inventoryServer = grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(inventoryServer, inventoryservice.NewService(stock))
	serveWG.Add(1)
	go func() {
		defer serveWG.Done()
		_ = inventoryServer.Serve(inventoryListener)
	}()
	deregisterInventory, err = registry.Register(discovery.Instance{
		Service: "inventory", ID: "inventory-demo", Address: inventoryListener.Addr().String(),
	})
	if err != nil {
		return fmt.Errorf("register inventory: %w", err)
	}

	connections = gateway.NewConnections(gateway.DefaultDial)
	handler := gateway.NewHandler(configStore, registry, connections, gateway.NewLimiter(nil))
	httpListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen gateway: %w", err)
	}
	httpServer = &http.Server{Handler: handler, ReadHeaderTimeout: time.Second}
	serveWG.Add(1)
	go func() {
		defer serveWG.Done()
		_ = httpServer.Serve(httpListener)
	}()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+httpListener.Addr().String()+"/api/v1/products/book-1",
		nil,
	)
	if err != nil {
		return fmt.Errorf("create demonstration request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer teaching-token")
	request.Header.Set("X-Request-Key", "chapter-14-demo")
	response, err := (&http.Client{Timeout: 2 * time.Second}).Do(request)
	if err != nil {
		return fmt.Errorf("call gateway: %w", err)
	}
	body, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if readErr != nil {
		return fmt.Errorf("read gateway response: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close gateway response: %w", closeErr)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway returned %s: %s", response.Status, body)
	}
	if _, err := out.Write(body); err != nil {
		return fmt.Errorf("write demonstration output: %w", err)
	}
	return nil
}
