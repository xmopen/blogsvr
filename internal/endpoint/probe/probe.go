package probe

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/commonlib/pkg/errcode"
)

// Health kubernetes HTTP 探针检测
func Health(c *gin.Context) {
	fmt.Println("Kubernetes Health Probe.")
	c.JSON(http.StatusOK, errcode.Success("success"))
}
