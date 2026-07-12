package upstream

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ScriptedResponse struct {
	StatusCode int
	Delay      time.Duration
	Product    Product
	Body       string
}

type ScriptedHandler struct {
	mu        sync.Mutex
	responses []ScriptedResponse
	calls     int
}

func NewScriptedHandler(responses []ScriptedResponse) *ScriptedHandler {
	copied := append([]ScriptedResponse(nil), responses...)
	if len(copied) == 0 {
		copied = []ScriptedResponse{{StatusCode: http.StatusOK, Product: Product{SKU: "book-1", Name: "Go Resilience", PriceCents: 9900, Quantity: 10}}}
	}
	return &ScriptedHandler{responses: copied}
}

func (h *ScriptedHandler) Calls() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.calls
}

func (h *ScriptedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h.next()
	if response.Delay > 0 {
		timer := time.NewTimer(response.Delay)
		select {
		case <-timer.C:
		case <-r.Context().Done():
			timer.Stop()
			return
		}
	}

	status := response.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status >= 200 && status < 300 {
		product := response.Product
		if product.SKU == "" {
			product.SKU = strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
		}
		_ = json.NewEncoder(w).Encode(product)
		return
	}
	if response.Body != "" {
		_, _ = w.Write([]byte(response.Body))
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(status)})
}

func (h *ScriptedHandler) next() ScriptedResponse {
	h.mu.Lock()
	defer h.mu.Unlock()
	index := h.calls
	h.calls++
	if index >= len(h.responses) {
		index = len(h.responses) - 1
	}
	return h.responses[index]
}
