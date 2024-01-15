// Package index  首页API.
package index

import (
	"net/http"
	"sync"

	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"

	"github.com/xmopen/commonlib/pkg/errcode"

	"github.com/xmopen/blogsvr/internal/models/articlemod"

	"github.com/xmopen/blogsvr/internal/manager/articlemanager"

	"github.com/gin-gonic/gin"
)

const (
	// listTypeIndex index page list.
	listTypeIndex = 0
	// listTypeHot index page hot list
	listTypeHot = 1
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

// indexRequest index request.
type indexRequest struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
	Type   int `json:"type" form:"type"`
}

// IndexArticleList  article list info for index page.
func (a *API) IndexArticleList(c *gin.Context) {
	request := &indexRequest{}
	xlog := ginhelper.Log(c)
	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusOK, errcode.ErrorGetIndexArticleList)
		return
	}
	xlog.Infof("index list req:[%+v]", request)
	if request.Limit <= 0 {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	articleList, err := a.articleListInfo(request)
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorGetIndexArticleList)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(a.convertModel(articleList)))
}

// articleListInfo get article list.
func (a *API) articleListInfo(request *indexRequest) ([]*articlemod.Article, error) {
	articles, err := articlemanager.Manager().AllPublishedArticles()
	if err != nil {
		return nil, err
	}
	if request.Type == listTypeIndex {
		if len(articles) < request.Offset {
			return nil, nil
		}
		if len(articles) >= request.Limit {
			return articles[request.Offset : request.Offset+request.Limit], nil
		}
		return articles[request.Offset:], nil
	}
	return articles, nil
}

func (a *API) convertModel(article []*articlemod.Article) []*articlemod.Article {
	responseList := make([]*articlemod.Article, 0)
	for _, item := range article {
		responseList = append(responseList, &articlemod.Article{
			ID:     item.ID,
			Img:    item.Img,
			Title:  item.Title,
			Time:   item.Time,
			Author: item.Author,
			Type:   item.Type,
		})
	}
	return responseList
}
