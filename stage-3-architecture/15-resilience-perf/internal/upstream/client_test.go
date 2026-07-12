package upstream

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestHTTPClientSuccess(t *testing.T) {
	srv := httptest.NewServer(NewScriptedHandler([]ScriptedResponse{{
		StatusCode: http.StatusOK,
		Product:    Product{SKU: "book-1", Name: "Go Microservices", PriceCents: 9900, Quantity: 7},
	}}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	got, err := client.GetProduct(context.Background(), "book-1")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}
	if got.SKU != "book-1" || got.Name != "Go Microservices" || got.PriceCents != 9900 || got.Quantity != 7 {
		t.Fatalf("unexpected product: %#v", got)
	}
}

func TestHTTPClientMapsServerErrorAsRetryable(t *testing.T) {
	srv := httptest.NewServer(NewScriptedHandler([]ScriptedResponse{{
		StatusCode: http.StatusServiceUnavailable,
		Body:       `{"error":"temporarily unavailable"}`,
	}}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.GetProduct(context.Background(), "book-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsRetryable(err) {
		t.Fatalf("expected retryable error, got %T %v", err, err)
	}
	var upstreamErr Error
	if !errors.As(err, &upstreamErr) {
		t.Fatalf("expected upstream Error, got %T", err)
	}
	if upstreamErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", upstreamErr.StatusCode)
	}
}

func TestHTTPClientDoesNotRetryClientErrors(t *testing.T) {
	srv := httptest.NewServer(NewScriptedHandler([]ScriptedResponse{{
		StatusCode: http.StatusNotFound,
		Body:       `{"error":"missing"}`,
	}}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.GetProduct(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if IsRetryable(err) {
		t.Fatalf("not-found error must not be retryable: %v", err)
	}
	if !IsNotFound(err) {
		t.Fatalf("expected not-found classification: %v", err)
	}
}

func TestScriptedHandlerUsesLastResponseAfterScriptIsExhausted(t *testing.T) {
	handler := NewScriptedHandler([]ScriptedResponse{{
		StatusCode: http.StatusServiceUnavailable,
		Body:       `{"error":"first"}`,
	}, {
		StatusCode: http.StatusOK,
		Product:    Product{SKU: "book-1", Name: "Recovered", PriceCents: 1, Quantity: 2},
	}})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, _ = client.GetProduct(context.Background(), "book-1")
	_, err := client.GetProduct(context.Background(), "book-1")
	if err != nil {
		t.Fatalf("second call should recover: %v", err)
	}
	_, err = client.GetProduct(context.Background(), "book-1")
	if err != nil {
		t.Fatalf("third call should reuse final response: %v", err)
	}
	if got := handler.Calls(); got != 3 {
		t.Fatalf("calls = %d", got)
	}
}

func TestHTTPClientMarksTransportErrorsRetryable(t *testing.T) {
	transportErr := errors.New("connection reset")
	client := NewHTTPClient("http://upstream.local", &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, transportErr
	})})

	_, err := client.GetProduct(context.Background(), "book-1")
	if err == nil {
		t.Fatal("expected transport error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected original transport error to be wrapped, got %v", err)
	}
	if !IsRetryable(err) {
		t.Fatalf("transport error should be retryable, got %T %v", err, err)
	}
}

func TestHTTPClientMarksTransportTimeoutRetryableUnlessContextExpired(t *testing.T) {
	transportErr := &url.Error{Op: "Get", URL: "http://upstream.local", Err: timeoutError{}}
	client := NewHTTPClient("http://upstream.local", &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, transportErr
	})})

	_, err := client.GetProduct(context.Background(), "book-1")
	if err == nil {
		t.Fatal("expected timeout transport error")
	}
	if !IsRetryable(err) {
		t.Fatalf("transport timeout should be retryable while request context has budget, got %T %v", err, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.GetProduct(ctx, "book-1")
	if err == nil {
		t.Fatal("expected canceled context error")
	}
	if IsRetryable(err) {
		t.Fatalf("context cancellation should not be retryable, got %T %v", err, err)
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestHTTPClientTimeoutWithoutContextExpiryIsRetryable(t *testing.T) {
	srv := httptest.NewServer(NewScriptedHandler([]ScriptedResponse{{
		StatusCode: http.StatusOK,
		Delay:      50 * time.Millisecond,
		Product:    Product{SKU: "book-1", Name: "Slow", PriceCents: 1, Quantity: 1},
	}}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, &http.Client{Timeout: time.Nanosecond})
	ctx := context.Background()
	_, err := client.GetProduct(ctx, "book-1")
	if err == nil {
		t.Fatal("expected client timeout error")
	}
	if ctx.Err() != nil {
		t.Fatalf("test setup expected caller context to remain healthy, got %v", ctx.Err())
	}
	if !IsRetryable(err) {
		t.Fatalf("client timeout should be retryable when caller context has budget, got %T %v", err, err)
	}
}

func TestIsRetryableDoesNotRetryContextWrappedURLErrors(t *testing.T) {
	for _, err := range []error{
		&url.Error{Op: "Get", URL: "http://upstream.local", Err: context.Canceled},
		&url.Error{Op: "Get", URL: "http://upstream.local", Err: context.DeadlineExceeded},
	} {
		if IsRetryable(err) {
			t.Fatalf("context-wrapped url error should not be retryable: %T %v", err, err)
		}
	}
}

func TestHTTPClientTreatsNilContextAsBackground(t *testing.T) {
	srv := httptest.NewServer(NewScriptedHandler([]ScriptedResponse{{
		StatusCode: http.StatusOK,
		Product:    Product{SKU: "book-1", Name: "Go Microservices", PriceCents: 9900, Quantity: 7},
	}}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("GetProduct(nil, ...) panicked: %v", recovered)
		}
	}()
	if _, err := client.GetProduct(nil, "book-1"); err != nil {
		t.Fatalf("GetProduct with nil context returned error: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
