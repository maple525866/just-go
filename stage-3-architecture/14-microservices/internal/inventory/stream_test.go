package inventory

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
)

func TestWatchStockStreamsCurrentAndUpdatedStock(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	client := newInventoryClient(t, NewService(store))
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.WatchStock(ctx, &inventoryv1.WatchStockRequest{Sku: "book-1"})
	if err != nil {
		t.Fatal(err)
	}

	initial, err := stream.Recv()
	if err != nil {
		t.Fatal(err)
	}
	if initial.GetQuantity() != 10 || initial.GetVersion() != 1 {
		t.Fatalf("initial = %#v", initial)
	}
	if _, err := store.Adjust("book-1", -2); err != nil {
		t.Fatal(err)
	}
	updated, err := stream.Recv()
	if err != nil {
		t.Fatal(err)
	}
	if updated.GetQuantity() != 8 || updated.GetVersion() != 2 {
		t.Fatalf("updated = %#v", updated)
	}

	cancel()
	if _, err := stream.Recv(); status.Code(err) != codes.Canceled {
		t.Fatalf("cancel code = %v, err = %v", status.Code(err), err)
	}
}

func TestGetStockUsesGRPCTransport(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	client := newInventoryClient(t, NewService(store))
	got, err := client.GetStock(context.Background(), &inventoryv1.GetStockRequest{Sku: "book-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetSku() != "book-1" || got.GetQuantity() != 10 || got.GetVersion() != 1 {
		t.Fatalf("response = %#v", got)
	}
}

func TestWatchStockMapsInvalidAndMissingSKU(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	client := newInventoryClient(t, NewService(store))

	tests := []struct {
		name string
		sku  string
		want codes.Code
	}{
		{name: "blank", sku: " ", want: codes.InvalidArgument},
		{name: "missing", sku: "missing", want: codes.NotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.WatchStock(context.Background(), &inventoryv1.WatchStockRequest{Sku: tt.sku})
			if err != nil {
				t.Fatal(err)
			}
			_, err = stream.Recv()
			if got := status.Code(err); got != tt.want {
				t.Fatalf("code = %v, want %v, err = %v", got, tt.want, err)
			}
		})
	}
}

func TestSyncStockStreamsOneResponsePerAdjustment(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	client := newInventoryClient(t, NewService(store))
	stream, err := client.SyncStock(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, step := range []struct {
		delta    int64
		quantity int64
		version  uint64
	}{
		{delta: -2, quantity: 8, version: 2},
		{delta: 5, quantity: 13, version: 3},
	} {
		if err := stream.Send(&inventoryv1.SyncStockRequest{Sku: "book-1", Delta: step.delta}); err != nil {
			t.Fatal(err)
		}
		got, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if got.GetQuantity() != step.quantity || got.GetVersion() != step.version {
			t.Fatalf("response = %#v", got)
		}
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Recv(); !errors.Is(err, io.EOF) {
		t.Fatalf("final Recv error = %v, want EOF", err)
	}
}

func TestSyncStockRejectsInvalidAdjustment(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	client := newInventoryClient(t, NewService(store))
	stream, err := client.SyncStock(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&inventoryv1.SyncStockRequest{Sku: "book-1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Recv(); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code = %v, err = %v", status.Code(err), err)
	}
}

func newInventoryClient(t *testing.T, service inventoryv1.InventoryServiceServer) inventoryv1.InventoryServiceClient {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(server, service)
	go func() {
		_ = server.Serve(listener)
	}()

	conn, err := grpc.NewClient(
		"passthrough:///inventory",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
		server.Stop()
		_ = listener.Close()
	})
	return inventoryv1.NewInventoryServiceClient(conn)
}
