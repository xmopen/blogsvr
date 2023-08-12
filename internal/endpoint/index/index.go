// Package index  首页API.
package index

import (
	"net/http"
	"sync"

	"github.com/xmopen/blogsvr/internal/errcode"

	"github.com/xmopen/blogsvr/internal/manager/articlemanager"

	"github.com/gin-gonic/gin"
)

var (
	apiInstance *API
	apiOnce     sync.Once
)

// API  index api.
type API struct {
}

// New  初始化Index API.
func New() *API {
	if apiInstance == nil {
		apiOnce.Do(func() {
			apiInstance = &API{}
		})
	}
	return apiInstance
}

// IndexArticleList  首页List.
func (a *API) IndexArticleList(c *gin.Context) {
	articles, err := articlemanager.Manager().AllPublishedArticles()
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorGetIndexArticleList)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(articles))
}

// IndexHotArticle 首页热门文章.
func (a *API) IndexHotArticle(c *gin.Context) {

}
