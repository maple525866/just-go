package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"just-go/stage-2-business/11-observability/healthx"
	"just-go/stage-2-business/11-observability/metricsx"
	"just-go/stage-2-business/11-observability/tracex"
	"just-go/stage-2-business/capstone-2-blog-api/auth"
	"just-go/stage-2-business/capstone-2-blog-api/cache"
	"just-go/stage-2-business/capstone-2-blog-api/model"
	"just-go/stage-2-business/capstone-2-blog-api/store"
)

type API struct {
	store   *store.MemoryStore
	cache   *cache.ArticleCache
	tokens  *auth.TokenManager
	metrics *metricsx.Registry
	health  *healthx.Checker
}

func NewAPI(s *store.MemoryStore, c *cache.ArticleCache, t *auth.TokenManager) http.Handler {
	h := healthx.NewChecker()
	h.AddLiveness("process", healthx.OK)
	h.AddReadiness("store", healthx.OK)
	return &API{store: s, cache: c, tokens: t, metrics: metricsx.NewRegistry(), health: h}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.metrics.Counter("blog_http_requests_total", "Total blog API requests").Add(1)
	tid := r.Header.Get("X-Trace-ID")
	ctx := r.Context()
	if tid != "" {
		ctx = tracex.WithTraceID(ctx, tid)
	}
	ctx, span := tracex.StartSpan(ctx, r.Method+" "+r.URL.Path)
	w.Header().Set("X-Trace-ID", span.TraceID)
	r = r.WithContext(ctx)
	switch {
	case r.URL.Path == "/livez":
		writeJSON(w, 200, a.health.Liveness(ctx))
	case r.URL.Path == "/readyz":
		rep := a.health.Readiness(ctx)
		code := 200
		if !rep.OK {
			code = 503
		}
		writeJSON(w, code, rep)
	case r.URL.Path == "/metrics":
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(a.metrics.Exposition()))
	case r.URL.Path == "/api/register" && r.Method == http.MethodPost:
		a.register(w, r)
	case r.URL.Path == "/api/login" && r.Method == http.MethodPost:
		a.login(w, r)
	case r.URL.Path == "/api/articles" && r.Method == http.MethodPost:
		a.createArticle(w, r)
	case r.URL.Path == "/api/articles" && r.Method == http.MethodGet:
		a.listArticles(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/articles/"):
		a.articleSubroutes(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (a *API) register(w http.ResponseWriter, r *http.Request) {
	var req struct{ Username, Password string }
	if !decode(w, r, &req) {
		return
	}
	h, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, 500, err)
		return
	}
	u, err := a.store.CreateUser(req.Username, h)
	if err != nil {
		writeError(w, status(err), err)
		return
	}
	writeJSON(w, 201, u)
}
func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var req struct{ Username, Password string }
	if !decode(w, r, &req) {
		return
	}
	u, err := a.store.UserByUsername(req.Username)
	if err != nil || !auth.CheckPassword(u.PasswordHash, req.Password) {
		writeError(w, 401, errors.New("invalid credentials"))
		return
	}
	tok, err := a.tokens.Sign(u.ID, u.Username)
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]string{"token": tok})
}
func (a *API) createArticle(w http.ResponseWriter, r *http.Request) {
	cl, ok := a.requireAuth(w, r)
	if !ok {
		return
	}
	var req model.ArticleInput
	if !decode(w, r, &req) {
		return
	}
	req.AuthorID = cl.UserID
	art, err := a.store.CreateArticle(req)
	if err != nil {
		writeError(w, status(err), err)
		return
	}
	a.cache.Invalidate(art.ID)
	a.metrics.Counter("blog_articles_created_total", "Articles created").Add(1)
	writeJSON(w, 201, art)
}
func (a *API) listArticles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	ps, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	out, err := a.store.ListArticles(model.ArticleFilter{Tag: r.URL.Query().Get("tag"), Page: page, PageSize: ps})
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, out)
}
func (a *API) articleSubroutes(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/articles/")
	parts := strings.Split(rest, "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("invalid article id"))
		return
	}
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			a.getArticle(w, id)
		case http.MethodPut:
			a.updateArticle(w, r, id)
		case http.MethodDelete:
			a.deleteArticle(w, r, id)
		default:
			http.NotFound(w, r)
		}
		return
	}
	if len(parts) == 2 && parts[1] == "comments" && r.Method == http.MethodPost {
		a.addComment(w, r, id)
		return
	}
	http.NotFound(w, r)
}

func (a *API) getArticle(w http.ResponseWriter, id int64) {
	if art, ok := a.cache.Get(id); ok {
		writeJSON(w, http.StatusOK, art)
		return
	}
	art, err := a.store.GetArticle(id)
	if err != nil {
		writeError(w, status(err), err)
		return
	}
	a.cache.Set(art)
	writeJSON(w, http.StatusOK, art)
}

func (a *API) updateArticle(w http.ResponseWriter, r *http.Request, id int64) {
	cl, ok := a.requireAuth(w, r)
	if !ok {
		return
	}
	var req model.ArticleInput
	if !decode(w, r, &req) {
		return
	}
	req.AuthorID = cl.UserID
	art, err := a.store.UpdateArticle(id, req)
	if err != nil {
		writeError(w, status(err), err)
		return
	}
	a.cache.Invalidate(id)
	a.metrics.Counter("blog_articles_updated_total", "Articles updated").Add(1)
	writeJSON(w, http.StatusOK, art)
}

func (a *API) deleteArticle(w http.ResponseWriter, r *http.Request, id int64) {
	if _, ok := a.requireAuth(w, r); !ok {
		return
	}
	if err := a.store.DeleteArticle(id); err != nil {
		writeError(w, status(err), err)
		return
	}
	a.cache.Invalidate(id)
	a.metrics.Counter("blog_articles_deleted_total", "Articles deleted").Add(1)
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) addComment(w http.ResponseWriter, r *http.Request, id int64) {
	cl, ok := a.requireAuth(w, r)
	if !ok {
		return
	}
	var req struct {
		Body     string `json:"body"`
		ParentID int64  `json:"parent_id"`
	}
	if !decode(w, r, &req) {
		return
	}
	c, err := a.store.AddComment(id, req.ParentID, cl.UserID, req.Body)
	if err != nil {
		writeError(w, status(err), err)
		return
	}
	a.cache.Invalidate(id)
	writeJSON(w, http.StatusCreated, c)
}
func (a *API) requireAuth(w http.ResponseWriter, r *http.Request) (auth.Claims, bool) {
	cl, err := a.tokens.Verify(auth.ParseBearer(r.Header.Get("Authorization")))
	if err != nil {
		writeError(w, 401, err)
		return auth.Claims{}, false
	}
	return cl, true
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, 400, err)
		return false
	}
	return true
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}
func status(err error) int {
	if errors.Is(err, store.ErrNotFound) {
		return 404
	}
	if errors.Is(err, store.ErrDuplicate) {
		return 409
	}
	if errors.Is(err, store.ErrInvalid) {
		return 422
	}
	return 500
}
func itoa(v int) string { return strconv.Itoa(v) }
