// Package endpoint API.
package endpoint

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/endpoint/article"
	"github.com/xmopen/blogsvr/internal/endpoint/index"
	"github.com/xmopen/blogsvr/internal/endpoint/middleware"
	"github.com/xmopen/blogsvr/internal/server/authsvr"
	"github.com/xmopen/commonlib/pkg/server/authserver"
)

// Init 初始化API.
func Init(r *gin.Engine) {
	middleware.InitMiddleWare(r)

	group := r.Group("/openxm/api/v1/index")
	indexAPI := index.New()
	group.GET("/list", indexAPI.IndexArticleList)
	group.GET("/host", indexAPI.IndexHotArticle)

	group = r.Group("/openxm/api/v1/article")
	articleAPI := article.New()
	group.GET("/info", articleAPI.ArticleInfo)

	r.GET("/goprc/test", func(c *gin.Context) {
		response := &authserver.AuthSvrResponse{}
		err := authsvr.Server().GetUserInfoByAccount(context.TODO(), &authserver.AuthSvrRequest{}, response)
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}
		d, _ := json.Marshal(response)
		c.JSON(http.StatusOK, string(d))
	})

}
