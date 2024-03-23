// Package articlemanager 文章管理器.
package articlemanager

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/blogsvr/internal/models/archivemod"
	"github.com/xmopen/blogsvr/internal/models/articlemod"
	"github.com/xmopen/commonlib/pkg/database/xmarchive"
	"github.com/xmopen/commonlib/pkg/protocol/xmeventprotocol"
	"github.com/xmopen/commonlib/pkg/xredis"
	"github.com/xmopen/golib/pkg/localcache/lru"
	"github.com/xmopen/golib/pkg/utils/timeutils"
	"github.com/xmopen/golib/pkg/xlogging"
)

var (
	articleManagerInstance     *ArticleManager
	articleManagerInstanceOnce sync.Once

	// xlog for article manager module.
	xlog = xlogging.Tag("article.manager")
)

// ArticleManager 文章管理器.
type ArticleManager struct {
	articleCache    *lru.LocalCache
	hotArticleCache *lru.LocalCache
}

// Manager 返回文章管理器. articlemanager.Manager()
func Manager() *ArticleManager {
	if articleManagerInstance == nil {
		articleManagerInstanceOnce.Do(func() {
			articleManagerInstance = &ArticleManager{
				articleCache:    lru.New(24*time.Hour, 128, loadAllPublishedArticles),
				hotArticleCache: lru.New(24*time.Hour, 15, loadHotArticleList),
			}
		})
		xredis.MultiSubScribe(config.BlogsRedis(), []string{string(xmeventprotocol.XMEventKeyOfArticleUpdate)},
			func(m *redis.Message) {
				fmt.Println("listener: " + m.String())
				articleManagerInstance.hotArticleCache.ClearAll()
				articleManagerInstance.articleCache.ClearAll()
			})
	}
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
	archiveMapping, err := getArchiveMapping()
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
		archive, ok := archiveMapping[article.TypeID]
		if ok {
			article.Type = archive.Name
		}
		articleTime, err := timeutils.StringTimeToCNTime(article.Time)
		if err != nil {
			xlog.Errorf("string time to cntime err:[%+v] source:[%+v]", err, article.Time)
			continue
		}
		article.Time = articleTime
	}
	publishedArticlesCacheValue := &articleCacheValue{
		allArticlesCache:  articles,
		articleID2Article: make(map[int]*articlemod.Article),
	}
	for _, item := range articles {
		publishedArticlesCacheValue.articleID2Article[item.ID] = item
	}
	return publishedArticlesCacheValue, nil
}

func loadHotArticleList(param any) (any, error) {
	articles, err := articlemod.AllArticles()
	if err != nil {
		return nil, err
	}
	archiveMapping, err := getArchiveMapping()
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
		// 不需要缓存Content
		article.Content = ""
		article.Author = ""
		article.SubHead = ""
		article.Img = ""
		archive, ok := archiveMapping[article.TypeID]
		if ok {
			article.Type = archive.Name
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].ReadCount >= articles[j].ReadCount
	})
	return articles, nil
}

func getArchiveMapping() (map[int]*xmarchive.XMBlogsArchive, error) {
	archiveList, err := archivemod.GetArchiveList()
	if err != nil {
		return nil, err
	}
	archiveMapping := make(map[int]*xmarchive.XMBlogsArchive)
	for _, item := range archiveList {
		archiveMapping[item.ID] = item
	}
	return archiveMapping, nil
}

// AllPublishedArticles 获取已经发布的所有文章.
func (a *ArticleManager) AllPublishedArticles() ([]*articlemod.Article, error) {
	itr, err := a.articleCache.Load("all_published_articles", "")
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	articleCache := itr.(*articleCacheValue)
	return articleCache.allArticlesCache, nil
}

// Article 通过ArticleID获取Article.
func (a *ArticleManager) Article(articleID int) (*articlemod.Article, error) {
	itr, err := a.articleCache.Load("all_published_articles", nil)
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	articleCache := itr.(*articleCacheValue)
	article := articleCache.articleID2Article[articleID]
	return article, nil
}

// GetHotArticleListWithLimit return host article list with limit
func (a *ArticleManager) GetHotArticleListWithLimit(limit int) ([]*articlemod.Article, error) {
	itr, err := a.hotArticleCache.Load("hot_articles", nil)
	if err != nil {
		return nil, err
	}
	articleList := itr.([]*articlemod.Article)
	if len(articleList) > limit {
		articleList = articleList[:limit]
	}
	return articleList, nil
}
