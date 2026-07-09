package cache

import (
	"sync"
	"time"

	"just-go/stage-2-business/capstone-2-blog-api/model"
)

type ArticleCache struct {
	mu    sync.Mutex
	ttl   time.Duration
	items map[int64]entry
}
type entry struct {
	article   model.Article
	expiresAt time.Time
}

func NewArticleCache(ttl time.Duration) *ArticleCache {
	return &ArticleCache{ttl: ttl, items: map[int64]entry{}}
}
func (c *ArticleCache) Get(id int64) (model.Article, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[id]
	if !ok || time.Now().After(e.expiresAt) {
		delete(c.items, id)
		return model.Article{}, false
	}
	return cloneArticle(e.article), true
}
func (c *ArticleCache) Set(a model.Article) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[a.ID] = entry{article: cloneArticle(a), expiresAt: time.Now().Add(c.ttl)}
}
func (c *ArticleCache) Invalidate(id int64) { c.mu.Lock(); defer c.mu.Unlock(); delete(c.items, id) }

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
