// Package endpoint API.
package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/endpoint/index"
)

// Init 初始化API.
func Init(r *gin.Engine) {
	group := r.Group("/openxm/api/v1/index")
	indexAPI := index.New()
	group.GET("/list", indexAPI.IndexArticleList)
	group.GET("/host", indexAPI.IndexHotArticle)
}
