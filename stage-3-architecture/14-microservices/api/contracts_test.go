package api_test

import (
	"testing"

	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
)

func TestGeneratedServiceDescriptors(t *testing.T) {
	if got := productv1.ProductService_ServiceDesc.ServiceName; got != "product.v1.ProductService" {
		t.Fatalf("product service name = %q", got)
	}
	if got := inventoryv1.InventoryService_ServiceDesc.ServiceName; got != "inventory.v1.InventoryService" {
		t.Fatalf("inventory service name = %q", got)
	}

	methods := inventoryv1.InventoryService_ServiceDesc.Methods
	if len(methods) != 1 || methods[0].MethodName != "GetStock" {
		t.Fatalf("unary methods = %#v", methods)
	}
	streams := inventoryv1.InventoryService_ServiceDesc.Streams
	if len(streams) != 2 {
		t.Fatalf("stream count = %d", len(streams))
	}
	if streams[0].StreamName != "WatchStock" || !streams[0].ServerStreams || streams[0].ClientStreams {
		t.Fatalf("WatchStock descriptor = %#v", streams[0])
	}
	if streams[1].StreamName != "SyncStock" || !streams[1].ServerStreams || !streams[1].ClientStreams {
		t.Fatalf("SyncStock descriptor = %#v", streams[1])
	}
}
