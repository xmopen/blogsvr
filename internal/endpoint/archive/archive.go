package archive

import (
	"net/http"
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
	request := &getArchiveRequest{}
	if err := c.ShouldBindQuery(request); err != nil {
		xlog.Errorf("bind query err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xlog.Infof("request:[%+v]", request)
	archiveList, err := archivemanager.Manager().GetArchiveList()
	if err != nil {
		xlog.Errorf("get archive list err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorSystemError)
		return
	}
	if request.Limit > 0 {
		if request.Limit > len(archiveList) {
			request.Limit = len(archiveList)
		}
		archiveList = archiveList[request.Offset:request.Limit]
	}
	if archiveList == nil {
		archiveList = make([]*xmarchive.XMBlogsArchive, 0)
	}
	c.JSON(http.StatusOK, errcode.Success(archiveList))
}
