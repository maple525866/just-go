package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"just-go/stage-3-architecture/12-clean-architecture/domain"
	"just-go/stage-3-architecture/12-clean-architecture/usecase"
)

type Handler struct {
	service *usecase.ArticleService
	mux     *http.ServeMux
}

type createArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type articleResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	Status      domain.Status `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	PublishedAt *time.Time    `json:"published_at,omitempty"`
}

func NewHandler(service *usecase.ArticleService) *Handler {
	h := &Handler{service: service, mux: http.NewServeMux()}
	h.mux.HandleFunc("POST /articles", h.create)
	h.mux.HandleFunc("GET /articles/{id}", h.get)
	h.mux.HandleFunc("POST /articles/{id}/publish", h.publish)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.mux.ServeHTTP(w, r) }

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request createArticleRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	input := usecase.CreateArticleInput{Title: request.Title, Body: request.Body}
	a, err := h.service.Create(r.Context(), input)
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toArticleResponse(a))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	a, err := h.service.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toArticleResponse(a))
}

func (h *Handler) publish(w http.ResponseWriter, r *http.Request) {
	a, err := h.service.Publish(r.Context(), r.PathValue("id"))
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toArticleResponse(a))
}

func toArticleResponse(a *domain.Article) articleResponse {
	return articleResponse{
		ID: a.ID, Title: a.Title, Body: a.Body, Status: a.Status,
		CreatedAt: a.CreatedAt, PublishedAt: a.PublishedAt,
	}
}

func writeApplicationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, usecase.ErrConflict), errors.Is(err, domain.ErrAlreadyPublished):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrIDRequired), errors.Is(err, domain.ErrTitleRequired), errors.Is(err, domain.ErrBodyRequired):
		writeError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
