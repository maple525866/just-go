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

func TestStoreRejectsForbiddenAndInvalidUpdates(t *testing.T) {
	s := NewMemoryStore()
	alice, _ := s.CreateUser("alice", "hash")
	bob, _ := s.CreateUser("bob", "hash")
	article, err := s.CreateArticle(model.ArticleInput{AuthorID: alice.ID, Title: "Original", Body: "body", Tags: []string{"go"}})
	if err != nil {
		t.Fatalf("CreateArticle: %v", err)
	}

	if _, err := s.UpdateArticle(article.ID, model.ArticleInput{AuthorID: bob.ID, Title: "Stolen"}); err != ErrForbidden {
		t.Fatalf("bob update err = %v, want ErrForbidden", err)
	}
	if err := s.DeleteArticle(article.ID, bob.ID); err != ErrForbidden {
		t.Fatalf("bob delete err = %v, want ErrForbidden", err)
	}

	invalid := []model.ArticleInput{
		{AuthorID: alice.ID},
		{AuthorID: alice.ID, Title: "   "},
		{AuthorID: alice.ID, Body: "   "},
	}
	for _, input := range invalid {
		if _, err := s.UpdateArticle(article.ID, input); err != ErrInvalid {
			t.Fatalf("UpdateArticle(%+v) err = %v, want ErrInvalid", input, err)
		}
	}

	updated, err := s.UpdateArticle(article.ID, model.ArticleInput{AuthorID: alice.ID, Tags: []string{"updated"}})
	if err != nil {
		t.Fatalf("valid update: %v", err)
	}
	if len(updated.Tags) != 1 || updated.Tags[0] != "updated" {
		t.Fatalf("updated tags = %+v", updated.Tags)
	}
	if err := s.DeleteArticle(article.ID, alice.ID); err != nil {
		t.Fatalf("alice delete: %v", err)
	}
}

func TestListArticlesHandlesHugePageWithoutPanic(t *testing.T) {
	s := NewMemoryStore()
	user, _ := s.CreateUser("alice", "hash")
	for _, title := range []string{"one", "two"} {
		if _, err := s.CreateArticle(model.ArticleInput{AuthorID: user.ID, Title: title, Body: "body"}); err != nil {
			t.Fatalf("CreateArticle: %v", err)
		}
	}

	page, err := s.ListArticles(model.ArticleFilter{Page: int(^uint(0) >> 1), PageSize: 100})
	if err != nil {
		t.Fatalf("ListArticles: %v", err)
	}
	if page.Total != 2 || len(page.Items) != 0 {
		t.Fatalf("page = %+v, want total 2 and no items", page)
	}
}
