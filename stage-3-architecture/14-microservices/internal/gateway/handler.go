package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
	"just-go/stage-3-architecture/14-microservices/internal/configcenter"
	"just-go/stage-3-architecture/14-microservices/internal/discovery"
)

const (
	productServiceName   = "product"
	inventoryServiceName = "inventory"
)

type ConfigReader interface {
	Current() (configcenter.Snapshot, error)
}

type Resolver interface {
	Resolve(service string) (discovery.Instance, error)
}

type ClientProvider interface {
	Product(context.Context, string) (productv1.ProductServiceClient, error)
	Inventory(context.Context, string) (inventoryv1.InventoryServiceClient, error)
}

type productDetails struct {
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	PriceCents   int64  `json:"price_cents"`
	Quantity     int64  `json:"quantity"`
	StockVersion uint64 `json:"stock_version"`
}

type Handler struct {
	config   ConfigReader
	resolver Resolver
	clients  ClientProvider
	limiter  *Limiter
}

func NewHandler(config ConfigReader, resolver Resolver, clients ClientProvider, limiter *Limiter) http.Handler {
	if limiter == nil {
		limiter = NewLimiter(nil)
	}
	handler := &Handler{config: config, resolver: resolver, clients: clients, limiter: limiter}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/products/{sku}", handler.getProductDetails)
	return mux
}

func (h *Handler) getProductDetails(w http.ResponseWriter, r *http.Request) {
	snapshot, err := h.config.Current()
	if err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "gateway configuration unavailable")
		return
	}
	config := snapshot.Config
	if !authorizedBearer(r.Header.Get("Authorization"), config.BearerToken) {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	requestKey := strings.TrimSpace(r.Header.Get("X-Request-Key"))
	if requestKey == "" {
		requestKey = r.RemoteAddr
	}
	if !h.limiter.Allow(requestKey, config.RateLimit, config.RateWindow) {
		writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}
	if !config.RouteEnabled || !configcenter.InRollout(requestKey, config.RolloutPercent) {
		writeJSONError(w, http.StatusNotFound, "route not found")
		return
	}
	sku := strings.TrimSpace(r.PathValue("sku"))
	if sku == "" {
		writeJSONError(w, http.StatusBadRequest, "sku is required")
		return
	}

	productInstance, err := h.resolver.Resolve(productServiceName)
	if err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "product service unavailable")
		return
	}
	inventoryInstance, err := h.resolver.Resolve(inventoryServiceName)
	if err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "inventory service unavailable")
		return
	}
	productClient, err := h.clients.Product(r.Context(), productInstance.Address)
	if err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "product service unavailable")
		return
	}
	inventoryClient, err := h.clients.Inventory(r.Context(), inventoryInstance.Address)
	if err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "inventory service unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.RequestTimeout)
	defer cancel()
	type downstreamResult struct {
		product   *productv1.GetProductResponse
		inventory *inventoryv1.GetStockResponse
		err       error
	}
	results := make(chan downstreamResult, 2)
	go func() {
		product, err := productClient.GetProduct(ctx, &productv1.GetProductRequest{Sku: sku})
		results <- downstreamResult{product: product, err: err}
	}()
	go func() {
		stock, err := inventoryClient.GetStock(ctx, &inventoryv1.GetStockRequest{Sku: sku})
		results <- downstreamResult{inventory: stock, err: err}
	}()

	var product *productv1.GetProductResponse
	var stock *inventoryv1.GetStockResponse
	for range 2 {
		result := <-results
		if result.err != nil {
			cancel()
			writeGatewayError(w, result.err)
			return
		}
		if result.product != nil {
			product = result.product
		}
		if result.inventory != nil {
			stock = result.inventory
		}
	}
	if product == nil || stock == nil {
		writeJSONError(w, http.StatusInternalServerError, "request failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(productDetails{
		SKU:          product.GetSku(),
		Name:         product.GetName(),
		PriceCents:   product.GetPriceCents(),
		Quantity:     stock.GetQuantity(),
		StockVersion: stock.GetVersion(),
	})
}
