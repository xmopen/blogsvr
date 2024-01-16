package probe

import (
	"net/http"

	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/commonlib/pkg/errcode"
)

// Health kubernetes HTTP 探针检测
func Health(c *gin.Context) {
	ginhelper.Log(c).Infof("kubernetes health probe")
	c.JSON(http.StatusOK, errcode.Success("success"))
}
