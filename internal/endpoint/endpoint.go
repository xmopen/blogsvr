// Package endpoint API.
package endpoint

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/endpoint/archive"
	"github.com/xmopen/blogsvr/internal/endpoint/article"
	"github.com/xmopen/blogsvr/internal/endpoint/comment"
	"github.com/xmopen/blogsvr/internal/endpoint/index"
	"github.com/xmopen/blogsvr/internal/endpoint/middleware"
	"github.com/xmopen/golib/pkg/xlogging"
)

// Init 初始化API.
func Init(r *gin.Engine) {
	middleware.InitMiddleWare(r)
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				xlogging.Tag("app.recover").Errorf("panic:[%+v]", string(debug.Stack()))
			}
		}()
	})

	group := r.Group("/openxm/api/v1/index")
	indexAPI := index.New()
	group.GET("/list", indexAPI.IndexArticleList)

	group = r.Group("/openxm/api/v1/article")
	articleAPI := article.New()
	group.GET("/info", articleAPI.ArticleInfo)
	group.GET("/hot", articleAPI.GetHotArticleList)

	commentAPI := comment.New()
	group = r.Group("/openxm/api/v1/comment")
	group.POST("/do", commentAPI.Comment)
	group.GET("/list", commentAPI.GetArticleCommentList)

	archiveAPI := archive.New()
	group = r.Group("/openxm/api/v1/archive")
	group.GET("/list", archiveAPI.GetArchiveList)
	group.GET("/list/article", archiveAPI.GetArchiveArticleList)
}
