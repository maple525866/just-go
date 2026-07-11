package product

import (
	"errors"
	"sync"
	"testing"
)

func TestCatalogGetReturnsNormalizedProduct(t *testing.T) {
	catalog, err := NewCatalog([]Product{{SKU: " book-1 ", Name: " Go Book ", PriceCents: 9900}})
	if err != nil {
		t.Fatal(err)
	}

	got, err := catalog.Get("book-1")
	if err != nil {
		t.Fatal(err)
	}
	want := Product{SKU: "book-1", Name: "Go Book", PriceCents: 9900}
	if got != want {
		t.Fatalf("product = %#v, want %#v", got, want)
	}

	got.Name = "changed"
	again, err := catalog.Get("book-1")
	if err != nil {
		t.Fatal(err)
	}
	if again.Name != "Go Book" {
		t.Fatalf("stored product changed to %#v", again)
	}
}

func TestNewCatalogRejectsInvalidProducts(t *testing.T) {
	tests := []struct {
		name     string
		products []Product
	}{
		{name: "blank sku", products: []Product{{Name: "Go Book", PriceCents: 1}}},
		{name: "blank name", products: []Product{{SKU: "book-1", PriceCents: 1}}},
		{name: "zero price", products: []Product{{SKU: "book-1", Name: "Go Book"}}},
		{name: "negative price", products: []Product{{SKU: "book-1", Name: "Go Book", PriceCents: -1}}},
		{name: "duplicate sku", products: []Product{
			{SKU: "book-1", Name: "Go Book", PriceCents: 1},
			{SKU: " book-1 ", Name: "Other", PriceCents: 2},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewCatalog(tt.products); !errors.Is(err, ErrInvalidProduct) {
				t.Fatalf("error = %v, want ErrInvalidProduct", err)
			}
		})
	}
}

func TestCatalogGetRejectsUnknownOrBlankSKU(t *testing.T) {
	catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.Get(" "); !errors.Is(err, ErrInvalidProduct) {
		t.Fatalf("blank SKU error = %v", err)
	}
	if _, err := catalog.Get("missing"); !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("missing SKU error = %v", err)
	}
}

func TestCatalogSupportsConcurrentReads(t *testing.T) {
	catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 1}})
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := catalog.Get("book-1"); err != nil {
				t.Errorf("Get: %v", err)
			}
		}()
	}
	wg.Wait()
}
