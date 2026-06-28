package store

import (
	"strconv"
	"sync"

	"just-go/stage-2-business/08-web-foundations/model"
)

// MemoryStore keeps article data in memory for the Web foundations chapter.
type MemoryStore struct {
	mu       sync.RWMutex
	articles []model.Article
	nextID   int
}

// NewSeededMemoryStore returns deterministic data used by examples and tests.
func NewSeededMemoryStore() *MemoryStore {
	return &MemoryStore{
		articles: []model.Article{
			{ID: "1", Title: "HTTP handlers are functions", Body: "A handler receives a request and writes a response.", Tags: []string{"http", "handler"}},
			{ID: "2", Title: "Middleware wraps handlers", Body: "Middleware can add logging, recovery, CORS, and request context.", Tags: []string{"middleware", "context"}},
		},
		nextID: 3,
	}
}

// List returns a copy of all articles.
func (s *MemoryStore) List() []model.Article {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]model.Article, len(s.articles))
	for i, article := range s.articles {
		items[i] = cloneArticle(article)
	}
	return items
}

// Get returns one article by ID.
func (s *MemoryStore) Get(id string) (model.Article, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, article := range s.articles {
		if article.ID == id {
			return cloneArticle(article), true
		}
	}
	return model.Article{}, false
}

// Create appends a new article and returns the stored value.
func (s *MemoryStore) Create(title, body string, tags []string) model.Article {
	s.mu.Lock()
	defer s.mu.Unlock()

	article := model.Article{
		ID:    strconv.Itoa(s.nextID),
		Title: title,
		Body:  body,
		Tags:  append([]string(nil), tags...),
	}
	s.nextID++
	s.articles = append(s.articles, article)
	return cloneArticle(article)
}

func cloneArticle(article model.Article) model.Article {
	article.Tags = append([]string(nil), article.Tags...)
	return article
}
