package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"just-go/stage-3-architecture/12-clean-architecture/infrastructure/memory"
	"just-go/stage-3-architecture/12-clean-architecture/usecase"
)

type testClock struct{ now time.Time }

func (c testClock) Now() time.Time { return c.now }

type testIDs struct {
	mu   sync.Mutex
	next int
}

func (g *testIDs) NewID() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	return "id-" + string(rune('0'+g.next))
}

func newTestHandler() http.Handler {
	service := usecase.NewArticleService(memory.NewArticleRepository(), testClock{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &testIDs{})
	return NewHandler(service)
}

func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestArticleHTTPWorkflow(t *testing.T) {
	h := newTestHandler()
	created := request(t, h, http.MethodPost, "/articles", `{"title":"Ports","body":"Adapters"}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", created.Code, created.Body.String())
	}
	var a articleResponse
	if err := json.NewDecoder(created.Body).Decode(&a); err != nil {
		t.Fatal(err)
	}
	got := request(t, h, http.MethodGet, "/articles/"+a.ID, "")
	if got.Code != http.StatusOK {
		t.Fatalf("status=%d", got.Code)
	}
	published := request(t, h, http.MethodPost, "/articles/"+a.ID+"/publish", "")
	if published.Code != http.StatusOK || !bytes.Contains(published.Body.Bytes(), []byte(`"status":"published"`)) {
		t.Fatalf("status=%d body=%s", published.Code, published.Body.String())
	}
	again := request(t, h, http.MethodPost, "/articles/"+a.ID+"/publish", "")
	if again.Code != http.StatusConflict {
		t.Fatalf("status=%d", again.Code)
	}
}

func TestConcurrentPublishOnlyOneRequestSucceeds(t *testing.T) {
	h := newTestHandler()
	created := request(t, h, http.MethodPost, "/articles", `{"title":"Ports","body":"Adapters"}`)
	var article articleResponse
	if err := json.NewDecoder(created.Body).Decode(&article); err != nil {
		t.Fatal(err)
	}

	statuses := make(chan int, 2)
	start := make(chan struct{})
	for range 2 {
		go func() {
			<-start
			req := httptest.NewRequest(http.MethodPost, "/articles/"+article.ID+"/publish", nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			statuses <- rec.Code
		}()
	}
	close(start)

	var successes, conflicts int
	for range 2 {
		switch <-statuses {
		case http.StatusOK:
			successes++
		case http.StatusConflict:
			conflicts++
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}
}

func TestArticleHTTPErrors(t *testing.T) {
	h := newTestHandler()
	tests := []struct {
		name, method, path, body string
		status                   int
	}{
		{"malformed", http.MethodPost, "/articles", `{`, http.StatusBadRequest},
		{"trailing JSON", http.MethodPost, "/articles", `{"title":"a","body":"b"} {}`, http.StatusBadRequest},
		{"unknown field", http.MethodPost, "/articles", `{"title":"a","body":"b","extra":true}`, http.StatusBadRequest},
		{"validation", http.MethodPost, "/articles", `{}`, http.StatusUnprocessableEntity},
		{"missing get", http.MethodGet, "/articles/missing", "", http.StatusNotFound},
		{"missing publish", http.MethodPost, "/articles/missing/publish", "", http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := request(t, h, tt.method, tt.path, tt.body)
			if got.Code != tt.status {
				t.Fatalf("status=%d body=%s", got.Code, got.Body.String())
			}
		})
	}
}
