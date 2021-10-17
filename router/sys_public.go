package router

import (
	"github.com/gin-gonic/gin"
)

// public router
func InitPublicRouter(r *gin.RouterGroup) (R gin.IRoutes) {
	r.Group("/public")
	{

	}
	return r
}
