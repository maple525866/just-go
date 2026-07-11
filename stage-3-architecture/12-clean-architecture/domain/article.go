package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrIDRequired       = errors.New("article id is required")
	ErrTitleRequired    = errors.New("article title is required")
	ErrBodyRequired     = errors.New("article body is required")
	ErrAlreadyPublished = errors.New("article is already published")
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

type Article struct {
	ID          string
	Title       string
	Body        string
	Status      Status
	CreatedAt   time.Time
	PublishedAt *time.Time
}

func NewArticle(id, title, body string, now time.Time) (*Article, error) {
	id, title, body = strings.TrimSpace(id), strings.TrimSpace(title), strings.TrimSpace(body)
	if id == "" {
		return nil, ErrIDRequired
	}
	if title == "" {
		return nil, ErrTitleRequired
	}
	if body == "" {
		return nil, ErrBodyRequired
	}
	return &Article{ID: id, Title: title, Body: body, Status: StatusDraft, CreatedAt: now.UTC()}, nil
}

func (a *Article) Publish(now time.Time) error {
	if a.Status == StatusPublished {
		return ErrAlreadyPublished
	}
	at := now.UTC()
	a.Status, a.PublishedAt = StatusPublished, &at
	return nil
}

func (a *Article) Clone() *Article {
	clone := *a
	if a.PublishedAt != nil {
		at := *a.PublishedAt
		clone.PublishedAt = &at
	}
	return &clone
}
