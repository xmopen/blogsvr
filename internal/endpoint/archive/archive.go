package archive

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/manager/archivemanager"
	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"
	"github.com/xmopen/commonlib/pkg/database/xmarchive"
	"github.com/xmopen/commonlib/pkg/errcode"
)

var (
	apiInstance *API
	initAPIOnce sync.Once
)

type getArchiveRequest struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}

// API return an archive api
type API struct{}

// New return a single api instance
func New() *API {
	initAPIOnce.Do(func() {
		apiInstance = &API{}
	})
	return apiInstance
}

// GetArchiveList  return archive list
func (a *API) GetArchiveList(c *gin.Context) {
	xlog := ginhelper.Log(c)
	xlog.Infof("查询分类")
	archiveList, err := archivemanager.Manager().GetArchiveList()
	if err != nil {
		fmt.Println(err.Error())
		xlog.Errorf("get archive list err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorSystemError)
		return
	}
	if archiveList == nil {
		archiveList = make([]*xmarchive.XMBlogsArchive, 0)
	}
	c.JSON(http.StatusOK, errcode.Success(archiveList))
}

// GetArchiveArticleList 获取分类下文章列表
func (a *API) GetArchiveArticleList(c *gin.Context) {
	archiveIDStr := c.Query("archive_id")
	if archiveIDStr == "" {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	archiveID, err := strconv.Atoi(archiveIDStr)
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	articleList, err := archivemanager.Manager().GetArticleWithArchiveID(archiveID)
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorSystemError)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(articleList))
}
