package pathreport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/models/reportmod"
	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"
	"github.com/xmopen/commonlib/pkg/errcode"
)

// PathReport 访问路径上报。
func PathReport(c *gin.Context) {
	path := c.Request.URL.Path
	// Nginx转发的时候通过X-Real-Ip接受真实的客户端IP。
	ip := c.GetHeader("X-Real-Ip")
	err := reportmod.Visit(&reportmod.VisitReport{
		Path: path,
		IP:   ip,
	})
	c.JSON(http.StatusOK, errcode.Success(nil))
	if err != nil {
		ginhelper.Log(c).Error("report err:[%+v]", err)
		return
	}
}

// PathReportWithRequestPath 通过请求参数path进行上报。
func PathReportWithRequestPath(c *gin.Context) {
	ip := c.GetHeader("X-Real-Ip")
	reportPath := c.Query("report_path")
	if reportPath == "" {
		c.JSON(http.StatusOK, errcode.Success(nil))
		return
	}
	err := reportmod.Visit(&reportmod.VisitReport{
		Path: reportPath,
		IP:   ip,
	})
	if err != nil {
		ginhelper.Log(c).Errorf("report vistit err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorSystemError)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(nil))
}
