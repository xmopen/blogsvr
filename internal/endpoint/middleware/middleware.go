package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/golib/pkg/middleware"
)

// InitMiddleWare init the gin middleware.
func InitMiddleWare(r *gin.Engine) {
	// use cros middleware
	r.Use(middleware.Cors())
}
