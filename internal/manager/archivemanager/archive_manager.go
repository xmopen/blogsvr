package archivemanager

import (
	"sort"
	"sync"
	"time"

	"github.com/xmopen/blogsvr/internal/models/archivemod"
	"github.com/xmopen/commonlib/pkg/database/xmarchive"
	"github.com/xmopen/golib/pkg/localcache"
)

var (
	archiveManagerInstance *ArchiveManager
	initArchiveManagerOnce sync.Once
)

// ArchiveManager archive manager
type ArchiveManager struct {
	archiveLocalCache *localcache.LocalCache
}

// Manager return an archive manager single instance
func Manager() *ArchiveManager {
	initArchiveManagerOnce.Do(func() {
		archiveManagerInstance = &ArchiveManager{
			archiveLocalCache: localcache.New(loadArchiveList, 32, 30*time.Minute),
		}
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

// GetArchiveList return blogs archive list
func (m *ArchiveManager) GetArchiveList() ([]*xmarchive.XMBlogsArchive, error) {
	itr, err := m.archiveLocalCache.LoadOrCreate("archive_list", nil)
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	return itr.([]*xmarchive.XMBlogsArchive), nil
}
