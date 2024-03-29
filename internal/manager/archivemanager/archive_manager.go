package archivemanager

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/xmopen/golib/pkg/localcache/lru"

	"github.com/redis/go-redis/v9"
	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/protocol/xmeventprotocol"
	"github.com/xmopen/commonlib/pkg/xredis"

	"github.com/xmopen/blogsvr/internal/manager/articlemanager"
	"github.com/xmopen/blogsvr/internal/models/articlemod"

	"github.com/xmopen/blogsvr/internal/models/archivemod"
	"github.com/xmopen/commonlib/pkg/database/xmarchive"
)

var (
	archiveManagerInstance *ArchiveManager
	initArchiveManagerOnce sync.Once
)

// ArchiveManager archive manager
type ArchiveManager struct {
	archiveLocalCache  *lru.LocalCache
	archiveArticleList *lru.LocalCache
}

// Manager return an archive manager single instance
func Manager() *ArchiveManager {
	initArchiveManagerOnce.Do(func() {
		archiveManagerInstance = &ArchiveManager{
			archiveLocalCache:  lru.New(24*time.Hour, 32, loadArchiveList),
			archiveArticleList: lru.New(24*time.Hour, 16, loadArchiveArticleList),
		}

		xredis.MultiSubScribe(config.BlogsRedis(), []string{string(xmeventprotocol.XMEventKeyOfArchiveUpdate)},
			func(m *redis.Message) {
				time.Sleep(5 * time.Second)
				fmt.Println("listener msg:" + m.String())
				archiveManagerInstance.archiveLocalCache.ClearAll()
				archiveManagerInstance.archiveArticleList.ClearAll()
			})
	})
	return archiveManagerInstance
}

func loadArchiveList(param any) (any, error) {
	archiveList, err := archivemod.GetArchiveList()
	if err != nil {
		return nil, err
	}
	if archiveList == nil {
		return nil, nil
	}
	sort.Slice(archiveList, func(i, j int) bool {
		return archiveList[i].Weight >= archiveList[j].Weight
	})
	return archiveList, nil
}

// loadArchiveArticleList 加载分类下所有的文章
func loadArchiveArticleList(param any) (any, error) {
	allArticles, err := articlemanager.Manager().AllPublishedArticles()
	if err != nil {
		return nil, err
	}
	if allArticles == nil {
		return nil, nil
	}
	archiveArticleListCache := make(map[int][]*articlemod.Article)
	for _, item := range allArticles {
		if _, ok := archiveArticleListCache[item.TypeID]; !ok {
			archiveArticleListCache[item.TypeID] = make([]*articlemod.Article, 0)
		}
		archiveArticleListCache[item.TypeID] = append(archiveArticleListCache[item.TypeID], item)
	}
	return archiveArticleListCache, nil
}

// GetArchiveList return blogs archive list
func (m *ArchiveManager) GetArchiveList() ([]*xmarchive.XMBlogsArchive, error) {
	itr, err := m.archiveLocalCache.Load("archive_list", nil)
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	return itr.([]*xmarchive.XMBlogsArchive), nil
}

func (m *ArchiveManager) GetArticleWithArchiveID(archiveID int) ([]*articlemod.Article, error) {
	itr, err := m.archiveArticleList.Load("archive_article_cache", archiveID)
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	archiveArticleCache := itr.(map[int][]*articlemod.Article)
	return archiveArticleCache[archiveID], nil
}
