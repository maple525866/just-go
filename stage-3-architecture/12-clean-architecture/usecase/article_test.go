package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"just-go/stage-3-architecture/12-clean-architecture/domain"
)

type fakeClock struct{ now time.Time }

func (f fakeClock) Now() time.Time { return f.now }

type fakeIDs struct{ id string }

func (f fakeIDs) NewID() string { return f.id }

type mockRepository struct {
	article                  *domain.Article
	createCalls, updateCalls int
	getErr                   error
}

func (m *mockRepository) Create(_ context.Context, a *domain.Article) error {
	m.createCalls++
	m.article = a.Clone()
	return nil
}
func (m *mockRepository) Get(_ context.Context, _ string) (*domain.Article, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.article.Clone(), nil
}
func (m *mockRepository) Update(_ context.Context, a *domain.Article, _ domain.Status) error {
	m.updateCalls++
	m.article = a.Clone()
	return nil
}

func TestArticleServiceCreateAndPublishWithMock(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := &mockRepository{}
	service := NewArticleService(repo, fakeClock{now}, fakeIDs{"fixed-id"})
	a, err := service.Create(context.Background(), CreateArticleInput{Title: "Ports", Body: "Adapters"})
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != "fixed-id" || repo.createCalls != 1 {
		t.Fatalf("unexpected create: %+v calls=%d", a, repo.createCalls)
	}
	published, err := service.Publish(context.Background(), a.ID)
	if err != nil {
		t.Fatal(err)
	}
	if published.Status != domain.StatusPublished || repo.updateCalls != 1 {
		t.Fatalf("unexpected publish: %+v", published)
	}
}

func TestArticleServiceDoesNotHideRepositoryError(t *testing.T) {
	want := ErrNotFound
	service := NewArticleService(&mockRepository{getErr: want}, fakeClock{}, fakeIDs{})
	_, got := service.Publish(context.Background(), "missing")
	if !errors.Is(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestArticleServiceRejectsInvalidBeforeRepository(t *testing.T) {
	repo := &mockRepository{}
	service := NewArticleService(repo, fakeClock{}, fakeIDs{"id"})
	_, err := service.Create(context.Background(), CreateArticleInput{})
	if !errors.Is(err, domain.ErrTitleRequired) || repo.createCalls != 0 {
		t.Fatalf("err=%v calls=%d", err, repo.createCalls)
	}
}
