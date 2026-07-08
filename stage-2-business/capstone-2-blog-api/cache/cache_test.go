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
