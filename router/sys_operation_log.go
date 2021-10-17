package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitOperationLogRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router := GetCasbinRouter(r, jwtOptions, "/operation/log")
	{
		router.GET("/list", v1.GetOperationLogs)
		router.DELETE("/delete/batch", v1.BatchDeleteOperationLogByIds)
	}
	return r
}
