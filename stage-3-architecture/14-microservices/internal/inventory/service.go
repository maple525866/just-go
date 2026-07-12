package inventory

import (
	"context"
	"errors"
	"io"

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

func (s *Service) WatchStock(req *inventoryv1.WatchStockRequest, stream inventoryv1.InventoryService_WatchStockServer) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	updates, err := s.store.Watch(stream.Context(), req.GetSku())
	if err != nil {
		return mapStoreError(err)
	}
	for {
		select {
		case <-stream.Context().Done():
			return status.FromContextError(stream.Context().Err()).Err()
		case stock, ok := <-updates:
			if !ok {
				if err := stream.Context().Err(); err != nil {
					return status.FromContextError(err).Err()
				}
				return nil
			}
			if err := stream.Send(&inventoryv1.WatchStockResponse{
				Sku:      stock.SKU,
				Quantity: stock.Quantity,
				Version:  stock.Version,
			}); err != nil {
				return err
			}
		}
	}
}

func (s *Service) SyncStock(stream inventoryv1.InventoryService_SyncStockServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			if ctxErr := stream.Context().Err(); ctxErr != nil {
				return status.FromContextError(ctxErr).Err()
			}
			return err
		}
		if req == nil {
			return status.Error(codes.InvalidArgument, "request is required")
		}
		stock, err := s.store.Adjust(req.GetSku(), req.GetDelta())
		if err != nil {
			return mapStoreError(err)
		}
		if err := stream.Send(&inventoryv1.SyncStockResponse{
			Sku:      stock.SKU,
			Quantity: stock.Quantity,
			Version:  stock.Version,
		}); err != nil {
			return err
		}
	}
}

func mapStoreError(err error) error {
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()
	case errors.Is(err, ErrInvalidStock):
		return status.Error(codes.InvalidArgument, "invalid stock request")
	case errors.Is(err, ErrStockNotFound):
		return status.Error(codes.NotFound, "stock not found")
	default:
		return status.Error(codes.Internal, "inventory operation failed")
	}
}
