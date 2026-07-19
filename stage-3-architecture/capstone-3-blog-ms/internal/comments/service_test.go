package comments

import (
	"context"
	"testing"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
)

func TestNestedCommentAndSoftDelete(t *testing.T) {
	service := NewService()
	root, err := service.CreateComment(context.Background(), &blogv1.CreateCommentRequest{PostId: 1, AuthorId: 1, Body: "root"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreateComment(context.Background(), &blogv1.CreateCommentRequest{PostId: 1, ParentId: root.Id, AuthorId: 2, Body: "reply"}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.DeleteComment(context.Background(), &blogv1.DeleteCommentRequest{PostId: 1, CommentId: root.Id, AuthorId: 1}); err != nil {
		t.Fatal(err)
	}
	list, _ := service.ListComments(context.Background(), &blogv1.ListCommentsRequest{PostId: 1})
	if !list.Comments[0].Deleted || list.Comments[0].Body != "" || len(list.Comments[0].Replies) != 1 {
		t.Fatalf("comments = %#v", list.Comments)
	}
}
