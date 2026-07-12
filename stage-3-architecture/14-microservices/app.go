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
	return runWithListen(ctx, out, net.Listen)
}

type listenFunc func(network, address string) (net.Listener, error)

func runWithListen(ctx context.Context, out io.Writer, listen listenFunc) (retErr error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	if out == nil {
		return errors.New("output writer is required")
	}
	if listen == nil {
		return errors.New("listen function is required")
	}
	requestCtx, cancelRequest := context.WithCancel(ctx)
	defer cancelRequest()

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
		serveErrors         = make(chan error, 3)
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

	productListener, err := listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen product grpc: %w", err)
	}
	productServer = grpc.NewServer()
	productv1.RegisterProductServiceServer(productServer, productservice.NewService(catalog))
	serveWG.Add(1)
	go func() {
		defer serveWG.Done()
		if err := productServer.Serve(productListener); !isNormalServeError(err) {
			serveErrors <- fmt.Errorf("product grpc serve: %w", err)
		}
	}()
	deregisterProduct, err = registry.Register(discovery.Instance{
		Service: "product", ID: "product-demo", Address: productListener.Addr().String(),
	})
	if err != nil {
		return fmt.Errorf("register product: %w", err)
	}

	inventoryListener, err := listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen inventory grpc: %w", err)
	}
	inventoryServer = grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(inventoryServer, inventoryservice.NewService(stock))
	serveWG.Add(1)
	go func() {
		defer serveWG.Done()
		if err := inventoryServer.Serve(inventoryListener); !isNormalServeError(err) {
			serveErrors <- fmt.Errorf("inventory grpc serve: %w", err)
		}
	}()
	deregisterInventory, err = registry.Register(discovery.Instance{
		Service: "inventory", ID: "inventory-demo", Address: inventoryListener.Addr().String(),
	})
	if err != nil {
		return fmt.Errorf("register inventory: %w", err)
	}

	connections = gateway.NewConnections(gateway.DefaultDial)
	handler := gateway.NewHandler(configStore, registry, connections, gateway.NewLimiter(nil))
	httpListener, err := listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen gateway: %w", err)
	}
	httpServer = &http.Server{Handler: handler, ReadHeaderTimeout: time.Second}
	serveWG.Add(1)
	go func() {
		defer serveWG.Done()
		if err := httpServer.Serve(httpListener); !isNormalServeError(err) {
			serveErrors <- fmt.Errorf("gateway serve: %w", err)
		}
	}()

	demoResult := make(chan struct {
		body []byte
		err  error
	}, 1)
	go func() {
		body, err := callDemo(requestCtx, httpListener.Addr().String())
		demoResult <- struct {
			body []byte
			err  error
		}{body: body, err: err}
	}()

	select {
	case err := <-serveErrors:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case result := <-demoResult:
		if result.err != nil {
			return result.err
		}
		if _, err := out.Write(result.body); err != nil {
			return fmt.Errorf("write demonstration output: %w", err)
		}
		return nil
	}
}

func callDemo(ctx context.Context, address string) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+address+"/api/v1/products/book-1",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create demonstration request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer teaching-token")
	request.Header.Set("X-Request-Key", "chapter-14-demo")
	response, err := (&http.Client{Timeout: 2 * time.Second}).Do(request)
	if err != nil {
		return nil, fmt.Errorf("call gateway: %w", err)
	}
	body, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read gateway response: %w", readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close gateway response: %w", closeErr)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway returned %s: %s", response.Status, body)
	}
	return body, nil
}

func isNormalServeError(err error) bool {
	return err == nil || errors.Is(err, grpc.ErrServerStopped) || errors.Is(err, http.ErrServerClosed)
}
