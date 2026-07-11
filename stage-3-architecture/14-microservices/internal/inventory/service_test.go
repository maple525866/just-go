package inventory

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
)

func TestServiceGetStock(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewService(store).GetStock(context.Background(), &inventoryv1.GetStockRequest{Sku: "book-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetSku() != "book-1" || got.GetQuantity() != 10 || got.GetVersion() != 1 {
		t.Fatalf("response = %#v", got)
	}
}

func TestServiceGetStockMapsErrors(t *testing.T) {
	store, err := NewStore(map[string]int64{"book-1": 10})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(store)

	tests := []struct {
		name string
		ctx  context.Context
		req  *inventoryv1.GetStockRequest
		want codes.Code
	}{
		{name: "nil request", ctx: context.Background(), want: codes.InvalidArgument},
		{name: "blank sku", ctx: context.Background(), req: &inventoryv1.GetStockRequest{Sku: " "}, want: codes.InvalidArgument},
		{name: "missing stock", ctx: context.Background(), req: &inventoryv1.GetStockRequest{Sku: "missing"}, want: codes.NotFound},
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests = append(tests, struct {
		name string
		ctx  context.Context
		req  *inventoryv1.GetStockRequest
		want codes.Code
	}{name: "canceled", ctx: canceled, req: &inventoryv1.GetStockRequest{Sku: "book-1"}, want: codes.Canceled})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetStock(tt.ctx, tt.req)
			if got := status.Code(err); got != tt.want {
				t.Fatalf("code = %v, want %v, err = %v", got, tt.want, err)
			}
		})
	}
}
