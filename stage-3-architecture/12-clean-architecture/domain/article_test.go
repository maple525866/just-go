package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewArticleValidation(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.FixedZone("test", 3600))
	article, err := NewArticle(" a-1 ", " Clean Go ", " boundaries ", now)
	if err != nil {
		t.Fatal(err)
	}
	if article.ID != "a-1" || article.Title != "Clean Go" || article.Status != StatusDraft {
		t.Fatalf("unexpected article: %+v", article)
	}
	if article.CreatedAt.Location() != time.UTC {
		t.Fatalf("created time must be UTC: %v", article.CreatedAt)
	}

	tests := []struct {
		name, id, title, body string
		want                  error
	}{
		{"id", " ", "title", "body", ErrIDRequired},
		{"title", "1", " ", "body", ErrTitleRequired},
		{"body", "1", "title", " ", ErrBodyRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := NewArticle(tt.id, tt.title, tt.body, now)
			if !errors.Is(got, tt.want) {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestArticlePublish(t *testing.T) {
	a, _ := NewArticle("1", "title", "body", time.Time{})
	now := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	if err := a.Publish(now); err != nil {
		t.Fatal(err)
	}
	if a.Status != StatusPublished || a.PublishedAt == nil || !a.PublishedAt.Equal(now) {
		t.Fatalf("unexpected published article: %+v", a)
	}
	if err := a.Publish(now); !errors.Is(err, ErrAlreadyPublished) {
		t.Fatalf("got %v", err)
	}
}

func TestCloneIsIndependent(t *testing.T) {
	a, _ := NewArticle("1", "title", "body", time.Time{})
	_ = a.Publish(time.Now())
	clone := a.Clone()
	clone.Title = "changed"
	*clone.PublishedAt = time.Time{}
	if a.Title == clone.Title || a.PublishedAt.IsZero() {
		t.Fatal("clone shares mutable state")
	}
}
