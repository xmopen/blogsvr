// Package articlemanager 文章管理器.
package articlemanager

import (
	"sync"
	"time"

	"github.com/xmopen/blogsvr/internal/models/articlemod"

	"github.com/xmopen/golib/pkg/localcache"
)

var (
	articleManagerInstance     *ArticleManager
	articleManagerInstanceOnce sync.Once
)

// ArticleManager 文章管理器.
type ArticleManager struct {
	articleCache *localcache.LocalCache
}

// Manager 返回文章管理器. articlemanager.Manager()
func Manager() *ArticleManager {
	if articleManagerInstance == nil {
		articleManagerInstanceOnce.Do(func() {
			articleManagerInstance = &ArticleManager{
				articleCache: localcache.New(loadAllPublishedArticles, 128, 1*time.Hour),
			}
		})
	}
	// TODO: 后续是否进行PUBLISH通知.
	return articleManagerInstance
}

// articleCacheValue 文章缓存结构.
type articleCacheValue struct {
	allArticlesCache  []*articlemod.Article
	articleID2Article map[int]*articlemod.Article
}

// loadAllPublishedArticles 加载所有已发布文章到内存.
func loadAllPublishedArticles(param any) (any, error) {
	articles, err := articlemod.AllArticles()
	if err != nil {
		return nil, err
	}
	cacheValue := &articleCacheValue{
		allArticlesCache:  articles,
		articleID2Article: make(map[int]*articlemod.Article),
	}
	for _, item := range articles {
		cacheValue.articleID2Article[item.ID] = item
	}
	return cacheValue, nil
}

// AllPublishedArticles 获取已经发布的所有文章.
func (a *ArticleManager) AllPublishedArticles() ([]*articlemod.Article, error) {
	itr, err := a.articleCache.LoadOrCreate("all_published_articles", "")
	if err != nil {
		return nil, err
	}
	articleCache := itr.(*articleCacheValue)
	return articleCache.allArticlesCache, nil
}

// Article 通过ArticleID获取Article.
func (a *ArticleManager) Article(articleID int) (*articlemod.Article, error) {
	itr, err := a.articleCache.LoadOrCreate("all_published_articles", "")
	if err != nil {
		return nil, err
	}
	articleCache := itr.(*articleCacheValue)
	article := articleCache.articleID2Article[articleID]
	return article, nil
}
