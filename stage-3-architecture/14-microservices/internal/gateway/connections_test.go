package gateway

import (
	"context"
	"errors"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestConnectionsCachesByAddress(t *testing.T) {
	var mu sync.Mutex
	dials := make(map[string]int)
	connections := NewConnections(func(_ context.Context, address string) (*grpc.ClientConn, error) {
		mu.Lock()
		dials[address]++
		mu.Unlock()
		return grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	defer connections.Close()

	ctx := context.Background()
	if _, err := connections.Product(ctx, "127.0.0.1:5001"); err != nil {
		t.Fatal(err)
	}
	if _, err := connections.Inventory(ctx, "127.0.0.1:5001"); err != nil {
		t.Fatal(err)
	}
	if _, err := connections.Product(ctx, "127.0.0.1:5002"); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	if dials["127.0.0.1:5001"] != 1 || dials["127.0.0.1:5002"] != 1 {
		t.Fatalf("dial counts = %#v", dials)
	}
}

func TestConnectionsDoesNotCacheDialFailure(t *testing.T) {
	want := errors.New("dial failed")
	count := 0
	connections := NewConnections(func(context.Context, string) (*grpc.ClientConn, error) {
		count++
		return nil, want
	})
	defer connections.Close()

	for range 2 {
		if _, err := connections.Product(context.Background(), "127.0.0.1:5001"); !errors.Is(err, want) {
			t.Fatalf("error = %v", err)
		}
	}
	if count != 2 {
		t.Fatalf("dial count = %d, want 2", count)
	}
}

func TestConnectionsCloseIsIdempotentAndRejectsNewClients(t *testing.T) {
	connections := NewConnections(func(_ context.Context, address string) (*grpc.ClientConn, error) {
		return grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	if _, err := connections.Product(context.Background(), "127.0.0.1:5001"); err != nil {
		t.Fatal(err)
	}
	if err := connections.Close(); err != nil {
		t.Fatal(err)
	}
	if err := connections.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	if _, err := connections.Product(context.Background(), "127.0.0.1:5001"); !errors.Is(err, ErrConnectionsClosed) {
		t.Fatalf("Product after Close = %v", err)
	}
	if _, err := connections.Inventory(context.Background(), "127.0.0.1:5001"); !errors.Is(err, ErrConnectionsClosed) {
		t.Fatalf("Inventory after Close = %v", err)
	}
}

func TestConnectionsRejectsInvalidInputAndCanceledContext(t *testing.T) {
	connections := NewConnections(DefaultDial)
	defer connections.Close()
	if _, err := connections.Product(context.Background(), " "); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("blank address error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := connections.Product(ctx, "127.0.0.1:5001"); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled context error = %v", err)
	}
}
