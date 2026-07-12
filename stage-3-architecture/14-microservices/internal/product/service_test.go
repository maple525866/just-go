package product

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
)

func TestServiceGetProduct(t *testing.T) {
	catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 9900}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewService(catalog).GetProduct(context.Background(), &productv1.GetProductRequest{Sku: "book-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetSku() != "book-1" || got.GetName() != "Go Book" || got.GetPriceCents() != 9900 {
		t.Fatalf("response = %#v", got)
	}
}

func TestServiceGetProductMapsErrors(t *testing.T) {
	catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 9900}})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(catalog)

	tests := []struct {
		name string
		ctx  context.Context
		req  *productv1.GetProductRequest
		want codes.Code
	}{
		{name: "nil request", ctx: context.Background(), want: codes.InvalidArgument},
		{name: "blank sku", ctx: context.Background(), req: &productv1.GetProductRequest{Sku: " "}, want: codes.InvalidArgument},
		{name: "missing product", ctx: context.Background(), req: &productv1.GetProductRequest{Sku: "missing"}, want: codes.NotFound},
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests = append(tests, struct {
		name string
		ctx  context.Context
		req  *productv1.GetProductRequest
		want codes.Code
	}{name: "canceled", ctx: canceled, req: &productv1.GetProductRequest{Sku: "book-1"}, want: codes.Canceled})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetProduct(tt.ctx, tt.req)
			if got := status.Code(err); got != tt.want {
				t.Fatalf("code = %v, want %v, err = %v", got, tt.want, err)
			}
		})
	}
}
