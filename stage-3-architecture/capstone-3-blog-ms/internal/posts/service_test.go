package posts

import (
	"context"
	"testing"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
)

func TestPostLifecycleAndOwnership(t *testing.T) {
	service := NewService()
	post, err := service.CreatePost(context.Background(), &blogv1.CreatePostRequest{AuthorId: 1, Title: "DDD", Body: "bounded context", Tags: []string{"Go", "go"}})
	if err != nil || len(post.Tags) != 1 {
		t.Fatalf("post=%#v err=%v", post, err)
	}
	if _, err := service.UpdatePost(context.Background(), &blogv1.UpdatePostRequest{Id: post.Id, AuthorId: 2, Title: "bad"}); err == nil {
		t.Fatal("expected ownership error")
	}
}
