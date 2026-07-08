package store

import (
	"testing"

	"just-go/stage-2-business/capstone-2-blog-api/model"
)

func TestStoreArticlePaginationTagsAndComments(t *testing.T) {
	s := NewMemoryStore()
	user, err := s.CreateUser("alice", "hash")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if _, err := s.CreateUser("alice", "hash"); err == nil {
		t.Fatalf("duplicate user should fail")
	}

	a1, err := s.CreateArticle(model.ArticleInput{AuthorID: user.ID, Title: "Go API", Body: "body", Tags: []string{"go", "api"}})
	if err != nil {
		t.Fatalf("CreateArticle a1: %v", err)
	}
	_, _ = s.CreateArticle(model.ArticleInput{AuthorID: user.ID, Title: "Cache", Body: "body", Tags: []string{"go", "cache"}})
	page, err := s.ListArticles(model.ArticleFilter{Tag: "go", Page: 1, PageSize: 1})
	if err != nil {
		t.Fatalf("ListArticles: %v", err)
	}
	if page.Total != 2 || len(page.Items) != 1 || page.Page != 1 || page.PageSize != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}

	root, err := s.AddComment(a1.ID, 0, user.ID, "first")
	if err != nil {
		t.Fatalf("AddComment root: %v", err)
	}
	child, err := s.AddComment(a1.ID, root.ID, user.ID, "reply")
	if err != nil {
		t.Fatalf("AddComment child: %v", err)
	}
	if err := s.SoftDeleteComment(a1.ID, child.ID); err != nil {
		t.Fatalf("SoftDeleteComment: %v", err)
	}
	article, err := s.GetArticle(a1.ID)
	if err != nil {
		t.Fatalf("GetArticle: %v", err)
	}
	if len(article.Comments) != 1 || len(article.Comments[0].Replies) != 1 || !article.Comments[0].Replies[0].Deleted {
		t.Fatalf("nested soft-deleted comment missing: %+v", article.Comments)
	}
}
