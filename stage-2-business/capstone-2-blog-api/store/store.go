package store

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"just-go/stage-2-business/capstone-2-blog-api/model"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("duplicate")
var ErrInvalid = errors.New("invalid input")

type MemoryStore struct {
	mu            sync.RWMutex
	nextUserID    int64
	nextArticleID int64
	nextCommentID int64
	users         map[int64]model.User
	usersByName   map[string]int64
	articles      map[int64]model.Article
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{nextUserID: 1, nextArticleID: 1, nextCommentID: 1, users: map[int64]model.User{}, usersByName: map[string]int64{}, articles: map[int64]model.Article{}}
}

func (s *MemoryStore) CreateUser(username, hash string) (model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	username = strings.TrimSpace(username)
	if username == "" || hash == "" {
		return model.User{}, ErrInvalid
	}
	if _, ok := s.usersByName[username]; ok {
		return model.User{}, ErrDuplicate
	}
	u := model.User{ID: s.nextUserID, Username: username, PasswordHash: hash}
	s.nextUserID++
	s.users[u.ID] = u
	s.usersByName[username] = u.ID
	return u, nil
}

func (s *MemoryStore) UserByUsername(username string) (model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.usersByName[username]
	if !ok {
		return model.User{}, ErrNotFound
	}
	return s.users[id], nil
}

func (s *MemoryStore) CreateArticle(in model.ArticleInput) (model.Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Body) == "" {
		return model.Article{}, ErrInvalid
	}
	now := time.Now().UTC()
	a := model.Article{ID: s.nextArticleID, AuthorID: in.AuthorID, Title: in.Title, Body: in.Body, Tags: dedupe(in.Tags), CreatedAt: now, UpdatedAt: now}
	s.nextArticleID++
	s.articles[a.ID] = a
	return cloneArticle(a), nil
}

func (s *MemoryStore) GetArticle(id int64) (model.Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.articles[id]
	if !ok {
		return model.Article{}, ErrNotFound
	}
	return cloneArticle(a), nil
}

func (s *MemoryStore) UpdateArticle(id int64, in model.ArticleInput) (model.Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[id]
	if !ok {
		return model.Article{}, ErrNotFound
	}
	if in.Title != "" {
		a.Title = in.Title
	}
	if in.Body != "" {
		a.Body = in.Body
	}
	if in.Tags != nil {
		a.Tags = dedupe(in.Tags)
	}
	a.UpdatedAt = time.Now().UTC()
	s.articles[id] = a
	return cloneArticle(a), nil
}

func (s *MemoryStore) DeleteArticle(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[id]; !ok {
		return ErrNotFound
	}
	delete(s.articles, id)
	return nil
}

func (s *MemoryStore) ListArticles(f model.ArticleFilter) (model.Page[model.Article], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 || f.PageSize > 100 {
		f.PageSize = 20
	}
	items := make([]model.Article, 0, len(s.articles))
	for _, a := range s.articles {
		if f.Tag == "" || hasTag(a.Tags, f.Tag) {
			items = append(items, cloneArticle(a))
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID > items[j].ID })
	total := len(items)
	start := (f.Page - 1) * f.PageSize
	if start > total {
		start = total
	}
	end := start + f.PageSize
	if end > total {
		end = total
	}
	return model.Page[model.Article]{Items: items[start:end], Total: total, Page: f.Page, PageSize: f.PageSize}, nil
}

func (s *MemoryStore) AddComment(articleID, parentID, authorID int64, body string) (model.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[articleID]
	if !ok {
		return model.Comment{}, ErrNotFound
	}
	if strings.TrimSpace(body) == "" {
		return model.Comment{}, ErrInvalid
	}
	c := model.Comment{ID: s.nextCommentID, ArticleID: articleID, ParentID: parentID, AuthorID: authorID, Body: body, CreatedAt: time.Now().UTC()}
	s.nextCommentID++
	if parentID == 0 {
		a.Comments = append(a.Comments, c)
	} else if !addReply(&a.Comments, parentID, c) {
		return model.Comment{}, ErrNotFound
	}
	s.articles[articleID] = a
	return c, nil
}

func (s *MemoryStore) SoftDeleteComment(articleID, commentID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[articleID]
	if !ok {
		return ErrNotFound
	}
	if !softDelete(&a.Comments, commentID) {
		return ErrNotFound
	}
	s.articles[articleID] = a
	return nil
}

func addReply(comments *[]model.Comment, parentID int64, reply model.Comment) bool {
	for i := range *comments {
		if (*comments)[i].ID == parentID {
			(*comments)[i].Replies = append((*comments)[i].Replies, reply)
			return true
		}
		if addReply(&(*comments)[i].Replies, parentID, reply) {
			return true
		}
	}
	return false
}
func softDelete(comments *[]model.Comment, id int64) bool {
	for i := range *comments {
		if (*comments)[i].ID == id {
			(*comments)[i].Deleted = true
			(*comments)[i].Body = ""
			return true
		}
		if softDelete(&(*comments)[i].Replies, id) {
			return true
		}
	}
	return false
}
func dedupe(tags []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, t := range tags {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
func hasTag(tags []string, tag string) bool {
	tag = strings.ToLower(tag)
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
func cloneArticle(a model.Article) model.Article {
	a.Tags = append([]string(nil), a.Tags...)
	a.Comments = cloneComments(a.Comments)
	return a
}
func cloneComments(in []model.Comment) []model.Comment {
	out := append([]model.Comment(nil), in...)
	for i := range out {
		out[i].Replies = cloneComments(out[i].Replies)
	}
	return out
}
