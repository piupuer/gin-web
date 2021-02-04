package router

import (
	"github.com/gin-gonic/gin"
)

// 公共路由, 任何人可访问
func InitPublicRouter(r *gin.RouterGroup) (R gin.IRoutes) {
	r.Group("/public")
	{

	}
	return r
}
