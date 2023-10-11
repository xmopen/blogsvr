package article

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/xmopen/blogsvr/internal/models/articlemod"

	"github.com/xmopen/blogsvr/internal/config"

	"github.com/xmopen/golib/pkg/xgoroutine"
	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/commonlib/pkg/errcode"

	"github.com/xmopen/blogsvr/internal/manager/articlemanager"
	"github.com/xmopen/blogsvr/internal/util/apputils"

	"github.com/gin-gonic/gin"
)

const (
	articleReadCountLimitClientIPExpired time.Duration = time.Hour
)

// API article getArticleInfoReport api.
type API struct {
}

// New a api instance.
func New() *API {
	return &API{}
}

// ArticleInfo return article getArticleInfoReport by article id.
func (a *API) ArticleInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	articleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xlog := apputils.Log(c)
	articleInfo, err := articlemanager.Manager().Article(articleID)
	if err != nil {
		xlog.Errorf("get article getArticleInfoReport err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xgoroutine.SafeGoroutine(func() {
		a.updateArticleReadCount(xlog, articleID, c.ClientIP())
	})
	c.JSON(http.StatusOK, errcode.Success(articleInfo))
}

// updateArticleReadCount update article read count,an ip address can only be updated once an hour
func (a *API) updateArticleReadCount(xlog *xlogging.Entry, articleID int, clientIP string) {
	updatedKey := fmt.Sprintf("xm_update_read_count_%d_%s", articleID, clientIP)
	if a.isUpdateArticle(xlog, updatedKey) {
		xlog.Infof("client ip haved updated ip:[%+v] article_id:[%+v]", clientIP, articleID)
		return
	}
	if err := articlemod.UpdateArticleReadCount(articleID); err != nil {
		xlog.Errorf("update article readcount err:[%+v] article_id:[%+v]", err, articleID)
		return
	}
	config.BlogsRedis().Set(context.TODO(), updatedKey, time.Now().Unix(), articleReadCountLimitClientIPExpired)
}

func (a *API) isUpdateArticle(xlog *xlogging.Entry, updateKey string) bool {
	isUpdated, err := config.BlogsRedis().Exists(context.TODO(), updateKey).Result()
	if err != nil {
		xlog.Errorf("get xm update red count limit err:[%+v] update key:[%+v]", err, updateKey)
	}
	return isUpdated > 0
}
