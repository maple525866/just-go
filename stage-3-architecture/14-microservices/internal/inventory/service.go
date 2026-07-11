package inventory

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
)

type Service struct {
	inventoryv1.UnimplementedInventoryServiceServer
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetStock(ctx context.Context, req *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	stock, err := s.store.Get(req.GetSku())
	if err != nil {
		return nil, mapStoreError(err)
	}
	return &inventoryv1.GetStockResponse{
		Sku:      stock.SKU,
		Quantity: stock.Quantity,
		Version:  stock.Version,
	}, nil
}

func mapStoreError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidStock):
		return status.Error(codes.InvalidArgument, "invalid stock request")
	case errors.Is(err, ErrStockNotFound):
		return status.Error(codes.NotFound, "stock not found")
	default:
		return status.Error(codes.Internal, "inventory operation failed")
	}
}
