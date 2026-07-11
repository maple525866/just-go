package memory

import (
	"context"
	"sync"

	"just-go/stage-3-architecture/12-clean-architecture/domain"
	"just-go/stage-3-architecture/12-clean-architecture/usecase"
)

type ArticleRepository struct {
	mu       sync.RWMutex
	articles map[string]*domain.Article
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{articles: make(map[string]*domain.Article)}
}

func (r *ArticleRepository) Create(_ context.Context, a *domain.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.articles[a.ID]; exists {
		return usecase.ErrConflict
	}
	r.articles[a.ID] = a.Clone()
	return nil
}

func (r *ArticleRepository) Get(_ context.Context, id string) (*domain.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, exists := r.articles[id]
	if !exists {
		return nil, usecase.ErrNotFound
	}
	return a.Clone(), nil
}

func (r *ArticleRepository) Update(_ context.Context, a *domain.Article, expected domain.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, exists := r.articles[a.ID]
	if !exists {
		return usecase.ErrNotFound
	}
	if current.Status != expected {
		return usecase.ErrConflict
	}
	r.articles[a.ID] = a.Clone()
	return nil
}
