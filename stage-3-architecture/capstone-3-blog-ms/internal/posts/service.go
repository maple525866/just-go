package posts

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/rpcerr"
)

type Service struct {
	blogv1.UnimplementedPostServiceServer
	mu     sync.RWMutex
	nextID int64
	posts  map[int64]*blogv1.Post
}

func NewService() *Service {
	return &Service{nextID: 1, posts: map[int64]*blogv1.Post{}}
}

func (s *Service) CreatePost(_ context.Context, req *blogv1.CreatePostRequest) (*blogv1.Post, error) {
	if req.GetAuthorId() <= 0 || strings.TrimSpace(req.GetTitle()) == "" || strings.TrimSpace(req.GetBody()) == "" {
		return nil, rpcerr.ToStatus(rpcerr.ErrInvalid)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	post := &blogv1.Post{Id: s.nextID, AuthorId: req.GetAuthorId(), Title: strings.TrimSpace(req.GetTitle()), Body: req.GetBody(), Tags: normalizeTags(req.GetTags()), CreatedAt: now, UpdatedAt: now}
	s.nextID++
	s.posts[post.Id] = clone(post)
	return clone(post), nil
}

func (s *Service) GetPost(_ context.Context, req *blogv1.GetPostRequest) (*blogv1.Post, error) {
	s.mu.RLock()
	post, ok := s.posts[req.GetId()]
	s.mu.RUnlock()
	if !ok {
		return nil, rpcerr.ToStatus(rpcerr.ErrNotFound)
	}
	return clone(post), nil
}

func (s *Service) ListPosts(_ context.Context, req *blogv1.ListPostsRequest) (*blogv1.ListPostsResponse, error) {
	page, size := int(req.GetPage()), int(req.GetPageSize())
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	s.mu.RLock()
	items := make([]*blogv1.Post, 0, len(s.posts))
	for _, post := range s.posts {
		if req.GetTag() == "" || contains(post.Tags, strings.ToLower(req.GetTag())) {
			items = append(items, clone(post))
		}
	}
	s.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool { return items[i].Id > items[j].Id })
	total := len(items)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}
	return &blogv1.ListPostsResponse{Posts: items[start:end], Total: int32(total), Page: int32(page), PageSize: int32(size)}, nil
}

func (s *Service) UpdatePost(_ context.Context, req *blogv1.UpdatePostRequest) (*blogv1.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	post, ok := s.posts[req.GetId()]
	if !ok {
		return nil, rpcerr.ToStatus(rpcerr.ErrNotFound)
	}
	if post.AuthorId != req.GetAuthorId() {
		return nil, rpcerr.ToStatus(rpcerr.ErrForbidden)
	}
	if strings.TrimSpace(req.GetTitle()) != "" {
		post.Title = strings.TrimSpace(req.GetTitle())
	}
	if strings.TrimSpace(req.GetBody()) != "" {
		post.Body = req.GetBody()
	}
	if req.GetReplaceTags() {
		post.Tags = normalizeTags(req.GetTags())
	}
	post.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	return clone(post), nil
}

func (s *Service) DeletePost(_ context.Context, req *blogv1.DeletePostRequest) (*blogv1.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	post, ok := s.posts[req.GetId()]
	if !ok {
		return nil, rpcerr.ToStatus(rpcerr.ErrNotFound)
	}
	if post.AuthorId != req.GetAuthorId() {
		return nil, rpcerr.ToStatus(rpcerr.ErrForbidden)
	}
	delete(s.posts, req.GetId())
	return &blogv1.Empty{}, nil
}

func normalizeTags(tags []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" && !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}
	return result
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func clone(post *blogv1.Post) *blogv1.Post {
	return &blogv1.Post{
		Id: post.Id, AuthorId: post.AuthorId, Title: post.Title, Body: post.Body,
		Tags: append([]string(nil), post.Tags...), CreatedAt: post.CreatedAt, UpdatedAt: post.UpdatedAt,
	}
}
