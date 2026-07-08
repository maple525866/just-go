package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"just-go/stage-2-business/capstone-2-blog-api/auth"
	"just-go/stage-2-business/capstone-2-blog-api/cache"
	"just-go/stage-2-business/capstone-2-blog-api/store"
)

func TestBlogAPISmokeFlow(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))

	register := request(t, api, http.MethodPost, "/api/register", `{"username":"alice","password":"secret"}`, "")
	if register.Code != http.StatusCreated {
		t.Fatalf("register = %d %s", register.Code, register.Body.String())
	}
	login := request(t, api, http.MethodPost, "/api/login", `{"username":"alice","password":"secret"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login = %d %s", login.Code, login.Body.String())
	}
	var authResp map[string]string
	_ = json.Unmarshal(login.Body.Bytes(), &authResp)
	token := authResp["token"]
	if token == "" {
		t.Fatalf("login response missing token: %s", login.Body.String())
	}

	created := request(t, api, http.MethodPost, "/api/articles", `{"title":"Stage 2","body":"capstone","tags":["go","api"]}`, token)
	if created.Code != http.StatusCreated {
		t.Fatalf("create article = %d %s", created.Code, created.Body.String())
	}
	var article map[string]any
	_ = json.Unmarshal(created.Body.Bytes(), &article)
	id := int(article["id"].(float64))

	list := request(t, api, http.MethodGet, "/api/articles?tag=go&page=1&page_size=10", "", "")
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "Stage 2") {
		t.Fatalf("list = %d %s", list.Code, list.Body.String())
	}
	comment := request(t, api, http.MethodPost, "/api/articles/"+itoa(id)+"/comments", `{"body":"nice"}`, token)
	if comment.Code != http.StatusCreated {
		t.Fatalf("comment = %d %s", comment.Code, comment.Body.String())
	}

	updated := request(t, api, http.MethodPut, "/api/articles/"+itoa(id), `{"title":"Stage 2 updated","body":"capstone updated","tags":["go","updated"]}`, token)
	if updated.Code != http.StatusOK || !strings.Contains(updated.Body.String(), "Stage 2 updated") {
		t.Fatalf("update = %d %s", updated.Code, updated.Body.String())
	}
	detail := request(t, api, http.MethodGet, "/api/articles/"+itoa(id), "", "")
	if detail.Code != http.StatusOK || !strings.Contains(detail.Body.String(), "Stage 2 updated") {
		t.Fatalf("detail after update = %d %s", detail.Code, detail.Body.String())
	}

	deleted := request(t, api, http.MethodDelete, "/api/articles/"+itoa(id), "", token)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete = %d %s", deleted.Code, deleted.Body.String())
	}
	missing := request(t, api, http.MethodGet, "/api/articles/"+itoa(id), "", "")
	if missing.Code != http.StatusNotFound {
		t.Fatalf("detail after delete = %d %s", missing.Code, missing.Body.String())
	}

	metrics := request(t, api, http.MethodGet, "/metrics", "", "")
	if metrics.Code != http.StatusOK || !strings.Contains(metrics.Body.String(), "blog_http_requests_total") {
		t.Fatalf("metrics = %d %s", metrics.Code, metrics.Body.String())
	}
}

func TestArticleIDValidationAndAuth(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))

	badID := request(t, api, http.MethodGet, "/api/articles/not-a-number", "", "")
	if badID.Code != http.StatusBadRequest {
		t.Fatalf("bad id = %d %s", badID.Code, badID.Body.String())
	}
	unauthorized := request(t, api, http.MethodPost, "/api/articles", `{"title":"x","body":"y"}`, "")
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized create = %d %s", unauthorized.Code, unauthorized.Body.String())
	}
}

func request(t *testing.T, h http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}
