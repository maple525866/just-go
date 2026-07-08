package server

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	chaptermw "just-go/stage-2-business/08-web-foundations/middleware"
	"just-go/stage-2-business/08-web-foundations/model"
	"just-go/stage-2-business/08-web-foundations/response"
	"just-go/stage-2-business/08-web-foundations/store"
	"just-go/stage-2-business/08-web-foundations/validation"
)

// NewStdMux demonstrates that net/http alone is enough to build a service.
func NewStdMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	return mux
}

// NewRouter builds the chi-based REST demo API.
func NewRouter(articles *store.MemoryStore) http.Handler {
	validator := validation.New()
	router := chi.NewRouter()

	router.Use(chaptermw.Recover)
	router.Use(chaptermw.RequestID)
	router.Use(chaptermw.Logger(slog.Default()))
	router.Use(chaptermw.CORS("*"))
	router.Use(chaptermw.NewLimiter(1_000).Middleware)

	router.Get("/healthz", healthz)
	router.Route("/api/articles", func(r chi.Router) {
		r.Get("/", listArticles(articles))
		r.Post("/", createArticle(articles, validator))
		r.Get("/{id}", getArticle(articles))
	})

	return router
}

func healthz(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func listArticles(articles *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := articles.List()
		response.JSON(w, http.StatusOK, model.ArticleListResponse{Items: items, Total: len(items)})
	}
}

func createArticle(articles *store.MemoryStore, validator *validation.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()

		var request model.CreateArticleRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON", nil)
			return
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			response.Error(w, http.StatusBadRequest, "invalid_json", "request body must contain a single JSON object", nil)
			return
		}

		if fields := validator.ValidateCreateArticle(request); len(fields) > 0 {
			response.Error(w, http.StatusUnprocessableEntity, "validation_failed", "request validation failed", fields)
			return
		}

		article := articles.Create(request.Title, request.Body, request.Tags)
		w.Header().Set("Location", "/api/articles/"+article.ID)
		response.JSON(w, http.StatusCreated, article)
	}
}

func getArticle(articles *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		article, ok := articles.Get(id)
		if !ok {
			response.Error(w, http.StatusNotFound, "article_not_found", "article not found", nil)
			return
		}
		response.JSON(w, http.StatusOK, article)
	}
}
