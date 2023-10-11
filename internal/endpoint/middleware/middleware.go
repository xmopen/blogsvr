package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"
)

// InitMiddleWare init the gin middleware.
func InitMiddleWare(r *gin.Engine) {
	// use cros middleware
	r.Use(ginhelper.Cors())
}
