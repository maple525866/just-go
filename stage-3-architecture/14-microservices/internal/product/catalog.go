package product

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrInvalidProduct  = errors.New("invalid product")
	ErrProductNotFound = errors.New("product not found")
)

type Product struct {
	SKU        string
	Name       string
	PriceCents int64
}

type Catalog struct {
	mu       sync.RWMutex
	products map[string]Product
}

func NewCatalog(products []Product) (*Catalog, error) {
	catalog := &Catalog{products: make(map[string]Product, len(products))}
	for _, candidate := range products {
		product := Product{
			SKU:        strings.TrimSpace(candidate.SKU),
			Name:       strings.TrimSpace(candidate.Name),
			PriceCents: candidate.PriceCents,
		}
		if product.SKU == "" || product.Name == "" || product.PriceCents <= 0 {
			return nil, fmt.Errorf("%w: sku, name, and positive price are required", ErrInvalidProduct)
		}
		if _, exists := catalog.products[product.SKU]; exists {
			return nil, fmt.Errorf("%w: duplicate sku %q", ErrInvalidProduct, product.SKU)
		}
		catalog.products[product.SKU] = product
	}
	return catalog, nil
}

func (c *Catalog) Get(sku string) (Product, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return Product{}, fmt.Errorf("%w: sku is required", ErrInvalidProduct)
	}

	c.mu.RLock()
	product, ok := c.products[sku]
	c.mu.RUnlock()
	if !ok {
		return Product{}, fmt.Errorf("%w: %s", ErrProductNotFound, sku)
	}
	return product, nil
}
