package article

import (
	"net/http"
	"strconv"

	"github.com/xmopen/blogsvr/internal/manager/articlemanager"
	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/blogsvr/internal/models/articlemod"
	"github.com/xmopen/blogsvr/internal/util/apputils"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/errcode"
)

// API article info api.
type API struct {
}

// New a api instance.
func New() *API {
	return &API{}
}

// ArticleInfo return article info by article id.
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
	articleInfo, err := a.info(xlog, articleID)
	if err != nil {
		xlog.Errorf("get article info err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(articleInfo))
}

// info return article info.
func (a *API) info(xlog *xlogging.Entry, id int) (*articlemod.Article, error) {
	articleInfo, err := articlemanager.Manager().Article(id)
	if err != nil {
		xlog.Errorf("articlemanager get article err:[%+v] id:[%+v]", err, id)
		return nil, err
	}
	return articleInfo, nil
}
