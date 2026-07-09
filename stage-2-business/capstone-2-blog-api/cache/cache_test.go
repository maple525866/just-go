package cache

import (
	"testing"
	"time"

	"just-go/stage-2-business/capstone-2-blog-api/model"
)

func TestArticleCacheSetGetInvalidateAndExpire(t *testing.T) {
	c := NewArticleCache(time.Minute)
	article := model.Article{ID: 9, Title: "cached"}
	c.Set(article)
	got, ok := c.Get(9)
	if !ok || got.Title != "cached" {
		t.Fatalf("cache get = %+v %v", got, ok)
	}
	c.Invalidate(9)
	if _, ok := c.Get(9); ok {
		t.Fatalf("cache should be invalidated")
	}

	expired := NewArticleCache(time.Nanosecond)
	expired.Set(article)
	time.Sleep(time.Millisecond)
	if _, ok := expired.Get(9); ok {
		t.Fatalf("cache should expire")
	}
}

func TestArticleCacheDoesNotExposeMutableSlices(t *testing.T) {
	c := NewArticleCache(time.Minute)
	article := model.Article{
		ID:    1,
		Title: "cached",
		Tags:  []string{"go"},
		Comments: []model.Comment{{
			ID:      1,
			Body:    "root",
			Replies: []model.Comment{{ID: 2, Body: "reply"}},
		}},
	}
	c.Set(article)

	article.Tags[0] = "changed-before-get"
	article.Comments[0].Body = "changed-before-get"
	article.Comments[0].Replies[0].Body = "changed-before-get"

	got, ok := c.Get(1)
	if !ok {
		t.Fatal("cache miss")
	}
	if got.Tags[0] != "go" || got.Comments[0].Body != "root" || got.Comments[0].Replies[0].Body != "reply" {
		t.Fatalf("cached article aliases original input: %+v", got)
	}

	got.Tags[0] = "changed-after-get"
	got.Comments[0].Body = "changed-after-get"
	got.Comments[0].Replies[0].Body = "changed-after-get"
	again, ok := c.Get(1)
	if !ok {
		t.Fatal("cache miss after mutation")
	}
	if again.Tags[0] != "go" || again.Comments[0].Body != "root" || again.Comments[0].Replies[0].Body != "reply" {
		t.Fatalf("cache exposed mutable slices: %+v", again)
	}
}
