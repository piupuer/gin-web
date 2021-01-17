package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 操作日志路由
func InitOperationLogRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := GetCasbinRouter(r, authMiddleware, "/operation/log")
	{
		router.GET("/list", v1.GetOperationLogs)
		router.DELETE("/delete/batch", v1.BatchDeleteOperationLogByIds)
	}
	return r
}
