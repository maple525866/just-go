package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneratedCompositionRoot(t *testing.T) {
	handler := initializeHandler()
	if handler == nil {
		t.Fatal("generated injector returned nil")
	}
	req := httptest.NewRequest(http.MethodGet, "/articles/missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}
