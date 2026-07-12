package product

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
)

type Service struct {
	productv1.UnimplementedProductServiceServer
	catalog *Catalog
}

func NewService(catalog *Catalog) *Service {
	return &Service{catalog: catalog}
}

func (s *Service) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	product, err := s.catalog.Get(req.GetSku())
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidProduct):
			return nil, status.Error(codes.InvalidArgument, "sku is required")
		case errors.Is(err, ErrProductNotFound):
			return nil, status.Error(codes.NotFound, "product not found")
		default:
			return nil, status.Error(codes.Internal, "product lookup failed")
		}
	}
	return &productv1.GetProductResponse{
		Sku:        product.SKU,
		Name:       product.Name,
		PriceCents: product.PriceCents,
	}, nil
}
