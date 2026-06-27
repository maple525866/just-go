package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "default name", path: "/", wantStatus: http.StatusOK, wantBody: "hello, gopher"},
		{name: "query name", path: "/?name=Ada", wantStatus: http.StatusOK, wantBody: "hello, Ada"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			HelloHandler().ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus || rec.Body.String() != tt.wantBody {
				t.Fatalf("handler = (%d, %q), want (%d, %q)", rec.Code, rec.Body.String(), tt.wantStatus, tt.wantBody)
			}
		})
	}
}

func TestFetchText(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "fetches from local test server", path: "/?name=Bob", want: "hello, Bob"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(HelloHandler())
			defer server.Close()
			got, err := FetchText(server.Client(), server.URL+tt.path)
			if err != nil {
				t.Fatalf("FetchText() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("FetchText() = %q, want %q", got, tt.want)
			}
		})
	}
}
