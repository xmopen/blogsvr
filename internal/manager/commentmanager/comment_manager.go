package commentmanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/xmopen/commonlib/pkg/database/xmcomment"

	"github.com/xmopen/blogsvr/internal/models/commentmod"
	"github.com/xmopen/golib/pkg/localcache"
)

var (
	commentManagerInstance *CommentManager
	commentManagerOnce     sync.Once
)

// CommentManager comment manager
type CommentManager struct {
	// commentLocalCache 评论缓存
	// TODO: 需要考虑数据一致性问题
	commentLocalCache *localcache.LocalCache
}

// Manager return a comment manager
func Manager() *CommentManager {
	commentManagerOnce.Do(func() {
		commentManagerInstance = &CommentManager{
			commentLocalCache: localcache.New(loadCommentByArticleID, 128, 1*time.Hour),
		}
	})
	return commentManagerInstance
}

func loadCommentByArticleID(param any) (any, error) {
	articleID, ok := param.(int)
	if !ok {
		return nil, fmt.Errorf("load comment param err source:[%+v]", param)
	}
	return commentmod.GetCommentListByArticleID(articleID)
}

// GetArticleComment get article comment list
func (m *CommentManager) GetArticleComment(articleID int) ([]*xmcomment.XMComment, error) {
	itr, err := m.commentLocalCache.LoadOrCreate(fmt.Sprintf("article_comment_%d", articleID), articleID)
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	return itr.([]*xmcomment.XMComment), nil
}
