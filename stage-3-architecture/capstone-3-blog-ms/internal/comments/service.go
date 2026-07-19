package comments

import (
	"context"
	"strings"
	"sync"
	"time"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/rpcerr"
)

type Service struct {
	blogv1.UnimplementedCommentServiceServer
	mu       sync.RWMutex
	nextID   int64
	comments map[int64][]*blogv1.Comment
}

func NewService() *Service {
	return &Service{nextID: 1, comments: map[int64][]*blogv1.Comment{}}
}

func (s *Service) CreateComment(_ context.Context, req *blogv1.CreateCommentRequest) (*blogv1.Comment, error) {
	if req.GetPostId() <= 0 || req.GetAuthorId() <= 0 || strings.TrimSpace(req.GetBody()) == "" {
		return nil, rpcerr.ToStatus(rpcerr.ErrInvalid)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	comment := &blogv1.Comment{Id: s.nextID, PostId: req.GetPostId(), ParentId: req.GetParentId(), AuthorId: req.GetAuthorId(), Body: req.GetBody(), CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if req.GetParentId() == 0 {
		s.comments[req.GetPostId()] = append(s.comments[req.GetPostId()], comment)
	} else if !addReply(s.comments[req.GetPostId()], req.GetParentId(), comment) {
		return nil, rpcerr.ToStatus(rpcerr.ErrNotFound)
	}
	s.nextID++
	return clone(comment), nil
}

func (s *Service) ListComments(_ context.Context, req *blogv1.ListCommentsRequest) (*blogv1.ListCommentsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &blogv1.ListCommentsResponse{Comments: cloneList(s.comments[req.GetPostId()])}, nil
}

func (s *Service) DeleteComment(_ context.Context, req *blogv1.DeleteCommentRequest) (*blogv1.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	comment := find(s.comments[req.GetPostId()], req.GetCommentId())
	if comment == nil {
		return nil, rpcerr.ToStatus(rpcerr.ErrNotFound)
	}
	if comment.AuthorId != req.GetAuthorId() {
		return nil, rpcerr.ToStatus(rpcerr.ErrForbidden)
	}
	comment.Deleted, comment.Body = true, ""
	return &blogv1.Empty{}, nil
}

func addReply(items []*blogv1.Comment, parentID int64, reply *blogv1.Comment) bool {
	for _, item := range items {
		if item.Id == parentID {
			item.Replies = append(item.Replies, reply)
			return true
		}
		if addReply(item.Replies, parentID, reply) {
			return true
		}
	}
	return false
}

func find(items []*blogv1.Comment, id int64) *blogv1.Comment {
	for _, item := range items {
		if item.Id == id {
			return item
		}
		if found := find(item.Replies, id); found != nil {
			return found
		}
	}
	return nil
}

func clone(item *blogv1.Comment) *blogv1.Comment {
	return &blogv1.Comment{
		Id: item.Id, PostId: item.PostId, ParentId: item.ParentId, AuthorId: item.AuthorId,
		Body: item.Body, Deleted: item.Deleted, CreatedAt: item.CreatedAt, Replies: cloneList(item.Replies),
	}
}

func cloneList(items []*blogv1.Comment) []*blogv1.Comment {
	result := make([]*blogv1.Comment, 0, len(items))
	for _, item := range items {
		result = append(result, clone(item))
	}
	return result
}
