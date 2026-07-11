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

func TestRegisterRejectsBlankPassword(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))

	for _, body := range []string{`{"username":"bob"}`, `{"username":"bob","password":""}`, `{"username":"bob","password":"   "}`} {
		rec := request(t, api, http.MethodPost, "/api/register", body, "")
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("register %s = %d %s, want 422", body, rec.Code, rec.Body.String())
		}
	}
}

func TestMalformedAuthorizationHeaderIsRejected(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))
	token := registerAndLogin(t, api, "alice")

	req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBufferString(`{"title":"Raw","body":"token","tags":["go"]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("raw token auth = %d %s, want 401", rec.Code, rec.Body.String())
	}
}

func TestDecodeRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))
	token := registerAndLogin(t, api, "alice")

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		token  string
	}{
		{name: "register unknown field", method: http.MethodPost, path: "/api/register", body: `{"username":"bob","password":"secret","role":"admin"}`},
		{name: "register trailing value", method: http.MethodPost, path: "/api/register", body: `{"username":"bob","password":"secret"} true`},
		{name: "create unknown field", method: http.MethodPost, path: "/api/articles", body: `{"title":"Stage 2","body":"capstone","tags":["go"],"admin":true}`, token: token},
		{name: "create trailing value", method: http.MethodPost, path: "/api/articles", body: `{"title":"Stage 2","body":"capstone","tags":["go"]} true`, token: token},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := request(t, api, tt.method, tt.path, tt.body, tt.token)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("response = %d %s, want 400", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestOnlyArticleAuthorCanUpdateOrDelete(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))
	alice := registerAndLogin(t, api, "alice")
	bob := registerAndLogin(t, api, "bob")

	created := request(t, api, http.MethodPost, "/api/articles", `{"title":"Owned","body":"by alice","tags":["go"]}`, alice)
	if created.Code != http.StatusCreated {
		t.Fatalf("create = %d %s", created.Code, created.Body.String())
	}
	id := articleID(t, created)

	bobUpdate := request(t, api, http.MethodPut, "/api/articles/"+itoa(id), `{"title":"Stolen"}`, bob)
	if bobUpdate.Code != http.StatusForbidden {
		t.Fatalf("bob update = %d %s, want 403", bobUpdate.Code, bobUpdate.Body.String())
	}
	bobDelete := request(t, api, http.MethodDelete, "/api/articles/"+itoa(id), "", bob)
	if bobDelete.Code != http.StatusForbidden {
		t.Fatalf("bob delete = %d %s, want 403", bobDelete.Code, bobDelete.Body.String())
	}

	aliceUpdate := request(t, api, http.MethodPut, "/api/articles/"+itoa(id), `{"title":"Updated"}`, alice)
	if aliceUpdate.Code != http.StatusOK || !strings.Contains(aliceUpdate.Body.String(), "Updated") {
		t.Fatalf("alice update = %d %s", aliceUpdate.Code, aliceUpdate.Body.String())
	}
	aliceDelete := request(t, api, http.MethodDelete, "/api/articles/"+itoa(id), "", alice)
	if aliceDelete.Code != http.StatusNoContent {
		t.Fatalf("alice delete = %d %s", aliceDelete.Code, aliceDelete.Body.String())
	}
}

func TestUpdateArticleRejectsInvalidBodies(t *testing.T) {
	api := NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("secret")))
	token := registerAndLogin(t, api, "alice")
	created := request(t, api, http.MethodPost, "/api/articles", `{"title":"Original","body":"body","tags":["go"]}`, token)
	if created.Code != http.StatusCreated {
		t.Fatalf("create = %d %s", created.Code, created.Body.String())
	}
	id := articleID(t, created)

	for _, body := range []string{`{}`, `{"title":"   "}`, `{"body":"   "}`} {
		rec := request(t, api, http.MethodPut, "/api/articles/"+itoa(id), body, token)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("update %s = %d %s, want 422", body, rec.Code, rec.Body.String())
		}
	}

	valid := request(t, api, http.MethodPut, "/api/articles/"+itoa(id), `{"tags":["updated"]}`, token)
	if valid.Code != http.StatusOK {
		t.Fatalf("valid tags update = %d %s", valid.Code, valid.Body.String())
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

func registerAndLogin(t *testing.T, h http.Handler, username string) string {
	t.Helper()
	register := request(t, h, http.MethodPost, "/api/register", `{"username":"`+username+`","password":"secret"}`, "")
	if register.Code != http.StatusCreated {
		t.Fatalf("register %s = %d %s", username, register.Code, register.Body.String())
	}
	login := request(t, h, http.MethodPost, "/api/login", `{"username":"`+username+`","password":"secret"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login %s = %d %s", username, login.Code, login.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(login.Body.Bytes(), &body); err != nil {
		t.Fatalf("login response is not JSON: %v", err)
	}
	if body["token"] == "" {
		t.Fatalf("login response missing token: %s", login.Body.String())
	}
	return body["token"]
}

func articleID(t *testing.T, rec *httptest.ResponseRecorder) int {
	t.Helper()
	var article map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &article); err != nil {
		t.Fatalf("article response is not JSON: %v", err)
	}
	return int(article["id"].(float64))
}
