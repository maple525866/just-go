package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"just-go/stage-2-business/11-observability/healthx"
	"just-go/stage-2-business/11-observability/metricsx"
)

func TestRouterHealthMetricsAndWork(t *testing.T) {
	checker := healthx.NewChecker()
	checker.AddLiveness("process", healthx.OK)
	checker.AddReadiness("database", healthx.OK)
	reg := metricsx.NewRegistry()
	router := NewRouter(checker, reg)

	for _, path := range []string{"/livez", "/readyz"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"ok":true`) {
			t.Fatalf("%s response = %d %s", path, rec.Code, rec.Body.String())
		}
	}

	work := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/work", nil)
	req.Header.Set("X-Trace-ID", "trace-http")
	router.ServeHTTP(work, req)
	if work.Code != http.StatusOK || work.Header().Get("X-Trace-ID") != "trace-http" || !strings.Contains(work.Body.String(), "trace-http") {
		t.Fatalf("work response = %d headers=%v body=%s", work.Code, work.Header(), work.Body.String())
	}

	metrics := httptest.NewRecorder()
	router.ServeHTTP(metrics, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := metrics.Body.String()
	if !strings.Contains(body, "http_requests_total") || !strings.Contains(body, "work_units_total") {
		t.Fatalf("metrics missing counters:\n%s", body)
	}
}

func TestRouterReadinessFailure(t *testing.T) {
	checker := healthx.NewChecker()
	checker.AddLiveness("process", healthx.OK)
	checker.AddReadiness("database", func(ctx context.Context) error { return healthx.ErrNotReady("db down") })
	router := NewRouter(checker, metricsx.NewRegistry())

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	result := rec.Result()
	data, _ := io.ReadAll(result.Body)
	if rec.Code != http.StatusServiceUnavailable || !strings.Contains(string(data), "db down") {
		t.Fatalf("readyz response = %d %s", rec.Code, string(data))
	}
	if got := result.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q, want application/json", got)
	}
}
