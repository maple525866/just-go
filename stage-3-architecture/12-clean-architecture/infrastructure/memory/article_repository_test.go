package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"just-go/stage-3-architecture/12-clean-architecture/domain"
	"just-go/stage-3-architecture/12-clean-architecture/usecase"
)

func TestArticleRepositoryLifecycleAndIsolation(t *testing.T) {
	ctx := context.Background()
	repo := NewArticleRepository()
	a, _ := domain.NewArticle("1", "title", "body", time.Time{})
	if err := repo.Create(ctx, a); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, a); !errors.Is(err, usecase.ErrConflict) {
		t.Fatalf("got %v", err)
	}
	a.Title = "outside mutation"
	got, err := repo.Get(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "title" {
		t.Fatal("create did not copy input")
	}
	got.Title = "read mutation"
	again, _ := repo.Get(ctx, "1")
	if again.Title != "title" {
		t.Fatal("get did not copy output")
	}
	again.Title = "updated"
	if err := repo.Update(ctx, again, domain.StatusDraft); err != nil {
		t.Fatal(err)
	}
	updated, _ := repo.Get(ctx, "1")
	if updated.Title != "updated" {
		t.Fatalf("got %+v", updated)
	}
}

func TestArticleRepositoryMissing(t *testing.T) {
	repo := NewArticleRepository()
	ctx := context.Background()
	if _, err := repo.Get(ctx, "missing"); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	a, _ := domain.NewArticle("missing", "t", "b", time.Time{})
	if err := repo.Update(ctx, a, domain.StatusDraft); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestArticleRepositoryUpdateIsCompareAndSwap(t *testing.T) {
	ctx := context.Background()
	repo := NewArticleRepository()
	draft, _ := domain.NewArticle("1", "title", "body", time.Time{})
	if err := repo.Create(ctx, draft); err != nil {
		t.Fatal(err)
	}

	first, _ := repo.Get(ctx, "1")
	second, _ := repo.Get(ctx, "1")
	_ = first.Publish(time.Unix(1, 0))
	_ = second.Publish(time.Unix(2, 0))

	results := make(chan error, 2)
	start := make(chan struct{})
	for _, article := range []*domain.Article{first, second} {
		go func() {
			<-start
			results <- repo.Update(ctx, article, domain.StatusDraft)
		}()
	}
	close(start)

	var successes, conflicts int
	for range 2 {
		switch err := <-results; {
		case err == nil:
			successes++
		case errors.Is(err, usecase.ErrConflict):
			conflicts++
		default:
			t.Fatalf("unexpected update error: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}
}
