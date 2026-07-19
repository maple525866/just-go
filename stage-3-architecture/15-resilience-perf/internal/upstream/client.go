package upstream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	GetProduct(context.Context, string) (Product, error)
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string, client *http.Client) *HTTPClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPClient{baseURL: strings.TrimRight(baseURL, "/"), client: client}
}

func (c *HTTPClient) GetProduct(ctx context.Context, sku string) (product Product, retErr error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return Product{}, err
	}
	if strings.TrimSpace(sku) == "" {
		return Product{}, Error{StatusCode: http.StatusBadRequest, Message: "sku is required"}
	}
	endpoint := c.baseURL + "/api/v1/products/" + url.PathEscape(sku)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return Product{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Product{}, ctxErr
		}
		return Product{}, transportError{err: err}
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			retErr = errors.Join(retErr, fmt.Errorf("close upstream response: %w", err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		msg := body.Error
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return Product{}, Error{StatusCode: resp.StatusCode, Temporary: resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500, Message: msg}
	}

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return Product{}, fmt.Errorf("decode upstream product: %w", err)
	}
	return product, nil
}
