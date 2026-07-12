package product

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
)

func TestGetProductUsesGRPCTransport(t *testing.T) {
	catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 9900}})
	if err != nil {
		t.Fatal(err)
	}
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer()
	productv1.RegisterProductServiceServer(server, NewService(catalog))
	go func() {
		_ = server.Serve(listener)
	}()

	conn, err := grpc.NewClient(
		"passthrough:///product",
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

	got, err := productv1.NewProductServiceClient(conn).GetProduct(
		context.Background(),
		&productv1.GetProductRequest{Sku: "book-1"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.GetSku() != "book-1" || got.GetName() != "Go Book" || got.GetPriceCents() != 9900 {
		t.Fatalf("response = %#v", got)
	}
}
