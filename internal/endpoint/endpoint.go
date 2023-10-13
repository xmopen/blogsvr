// Package endpoint API.
package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/endpoint/archive"
	"github.com/xmopen/blogsvr/internal/endpoint/article"
	"github.com/xmopen/blogsvr/internal/endpoint/comment"
	"github.com/xmopen/blogsvr/internal/endpoint/index"
	"github.com/xmopen/blogsvr/internal/endpoint/middleware"
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

	commentAPI := comment.New()
	group = r.Group("/openxm/api/v1/comment")
	group.POST("/do", commentAPI.Comment)
	group.GET("/list", commentAPI.GetArticleCommentList)

	archiveAPI := archive.New()
	group = r.Group("/openxm/api/v1/archive")
	group.GET("/list", archiveAPI.GetArchiveList)
}
