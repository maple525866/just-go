package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"just-go/stage-2-business/08-web-foundations/model"
	"just-go/stage-2-business/08-web-foundations/store"
)

func TestNewStdMuxHealthz(t *testing.T) {
	handler := NewStdMux()
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("health response is not JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("health status = %q, want ok", body["status"])
	}
}

func TestNewRouterSuccessfulRequests(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		assert     func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "health",
			method:     http.MethodGet,
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "list articles",
			method:     http.MethodGet,
			path:       "/api/articles",
			wantStatus: http.StatusOK,
			assert: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var body model.ArticleListResponse
				decodeJSON(t, recorder, &body)
				if body.Total != 2 || len(body.Items) != 2 {
					t.Fatalf("list response = %+v, want two seeded articles", body)
				}
			},
		},
		{
			name:       "create article",
			method:     http.MethodPost,
			path:       "/api/articles",
			body:       `{"title":"Testing handlers","body":"httptest makes handlers easy to verify.","tags":["test","http"]}`,
			wantStatus: http.StatusCreated,
			assert: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if got := recorder.Header().Get("Location"); got != "/api/articles/3" {
					t.Fatalf("Location = %q, want /api/articles/3", got)
				}
				var body model.Article
				decodeJSON(t, recorder, &body)
				if body.ID != "3" || body.Title != "Testing handlers" {
					t.Fatalf("created article = %+v", body)
				}
			},
		},
		{
			name:       "get article",
			method:     http.MethodGet,
			path:       "/api/articles/1",
			wantStatus: http.StatusOK,
			assert: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var body model.Article
				decodeJSON(t, recorder, &body)
				if body.ID != "1" {
					t.Fatalf("article ID = %q, want 1", body.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRouter(store.NewSeededMemoryStore())
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				request.Header.Set("Content-Type", "application/json")
			}

			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("content type = %q, want application/json", got)
			}
			if tt.assert != nil {
				tt.assert(t, recorder)
			}
		})
	}
}

func TestNewRouterErrorRequests(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{name: "invalid JSON", method: http.MethodPost, path: "/api/articles", body: `{`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "trailing JSON garbage", method: http.MethodPost, path: "/api/articles", body: `{"title":"Testing handlers","body":"httptest makes handlers easy to verify.","tags":["test"]} true`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "validation failure", method: http.MethodPost, path: "/api/articles", body: `{"title":"","body":"short"}`, wantStatus: http.StatusUnprocessableEntity, wantCode: "validation_failed"},
		{name: "blank strings validation failure", method: http.MethodPost, path: "/api/articles", body: `{"title":"   ","body":"          ","tags":["go"]}`, wantStatus: http.StatusUnprocessableEntity, wantCode: "validation_failed"},
		{name: "not found", method: http.MethodGet, path: "/api/articles/missing", wantStatus: http.StatusNotFound, wantCode: "article_not_found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRouter(store.NewSeededMemoryStore())
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))

			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			var body model.ErrorResponse
			decodeJSON(t, recorder, &body)
			if body.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", body.Code, tt.wantCode)
			}
		})
	}
}

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.Unmarshal(recorder.Body.Bytes(), target); err != nil {
		t.Fatalf("response is not JSON: %v; body=%s", err, recorder.Body.String())
	}
}
