package commentmanager

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/xmopen/golib/pkg/localcache/lru"

	"github.com/redis/go-redis/v9"
	"github.com/xmopen/commonlib/pkg/protocol/xmeventprotocol"
	"github.com/xmopen/commonlib/pkg/xredis"

	"github.com/xmopen/golib/pkg/utils/timeutils"

	"github.com/xmopen/blogsvr/internal/config"

	"github.com/xmopen/golib/pkg/xgoroutine"

	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/blogsvr/internal/server/authsvr"
	"github.com/xmopen/commonlib/pkg/server/authserver"

	"github.com/xmopen/blogsvr/internal/models/commentmod"
)

var (
	commentManagerInstance *CommentManager
	commentManagerOnce     sync.Once
)

// CommentManager comment manager
type CommentManager struct {
	// commentLocalCache comment cache
	commentLocalCache *lru.LocalCache
}

// Manager return a comment manager
func Manager() *CommentManager {
	commentManagerOnce.Do(func() {
		commentManagerInstance = &CommentManager{
			commentLocalCache: lru.New(24*time.Hour, 128, loadCommentByArticleID),
		}
		xgoroutine.SafeGoroutine(func() {
			xredis.MultiSubScribe(config.BlogsRedis(), []string{string(xmeventprotocol.XMEventKeyOfArticleCommentUpdate)},
				func(m *redis.Message) {
					fmt.Println("listener: " + m.String())
					commentManagerInstance.commentLocalCache.ClearAll()
				})
		})

	})
	return commentManagerInstance
}

func loadCommentByArticleID(param any) (any, error) {
	articleID, ok := param.(int)
	if !ok {
		return nil, fmt.Errorf("load comment param err source:[%+v]", param)
	}
	xlog := xlogging.Tag("localcache.comment")
	commentList, err := commentmod.GetCommentListByArticleID(articleID)
	if err != nil {
		xlog.Errorf("get article comment list err:[%+v]", err)
		return nil, err
	}
	resultList := make([]*commentmod.Comment, 0)
	for _, cmt := range commentList {
		// rpc 请求关联User等相关信息.
		response := &authserver.AuthSvrResponse{}
		err = authsvr.Server().GetUserInfoByAccount(context.TODO(), &authserver.AuthSvrRequest{
			XMAccount: cmt.Account,
		}, response)
		if err != nil {
			xlog.Errorf("get userinfo rpc err:[%+v]", err)
			continue
		}
		commentTime, err := timeutils.StringTimeToCNTime(cmt.CreateTime.Format(time.DateTime))
		if err != nil {
			xlog.Errorf("transforme time to cn time err:[%+v]", err)
			continue
		}
		resultList = append(resultList, &commentmod.Comment{
			Username:    response.XMUserInfo.UserName,
			Icon:        response.XMUserInfo.UserIcon,
			CommentTime: commentTime,
			XMComment:   cmt,
		})
	}
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].XMComment.CreateTime.Unix() > resultList[j].XMComment.CreateTime.Unix()
	})
	return resultList, nil
}

// GetArticleComment get article comment list
func (m *CommentManager) GetArticleComment(articleID int) ([]*commentmod.Comment, error) {
	itr, err := m.commentLocalCache.Load(fmt.Sprintf("article_comment_%d", articleID), articleID)
	if err != nil {
		return nil, err
	}
	if itr == nil {
		return nil, nil
	}
	return itr.([]*commentmod.Comment), nil
}
