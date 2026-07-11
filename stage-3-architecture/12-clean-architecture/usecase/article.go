package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"just-go/stage-3-architecture/12-clean-architecture/domain"
)

var (
	ErrNotFound = errors.New("article not found")
	ErrConflict = errors.New("article state conflict")
)

type ArticleRepository interface {
	Create(context.Context, *domain.Article) error
	Get(context.Context, string) (*domain.Article, error)
	Update(context.Context, *domain.Article, domain.Status) error
}

type Clock interface{ Now() time.Time }
type IDGenerator interface{ NewID() string }

type SystemClock struct{}

func NewSystemClock() Clock        { return SystemClock{} }
func (SystemClock) Now() time.Time { return time.Now() }

type SequentialIDGenerator struct{ next atomic.Uint64 }

func NewSequentialIDGenerator() IDGenerator    { return &SequentialIDGenerator{} }
func (g *SequentialIDGenerator) NewID() string { return fmt.Sprintf("article-%d", g.next.Add(1)) }

type ArticleService struct {
	repo  ArticleRepository
	clock Clock
	ids   IDGenerator
}

func NewArticleService(repo ArticleRepository, clock Clock, ids IDGenerator) *ArticleService {
	return &ArticleService{repo: repo, clock: clock, ids: ids}
}

type CreateArticleInput struct {
	Title string
	Body  string
}

func (s *ArticleService) Create(ctx context.Context, in CreateArticleInput) (*domain.Article, error) {
	a, err := domain.NewArticle(s.ids.NewID(), in.Title, in.Body, s.clock.Now())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a.Clone(), nil
}

func (s *ArticleService) Get(ctx context.Context, id string) (*domain.Article, error) {
	return s.repo.Get(ctx, id)
}

func (s *ArticleService) Publish(ctx context.Context, id string) (*domain.Article, error) {
	a, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := a.Publish(s.clock.Now()); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, a, domain.StatusDraft); err != nil {
		return nil, err
	}
	return a.Clone(), nil
}
