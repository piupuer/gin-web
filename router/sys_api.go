package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "go-shipment-api/api/v1"
)

// 接口路由
func InitApiRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("api").Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/list", v1.GetApis)
		router.GET("/role/category/:roleId", v1.GetRoleCategoryApisByRoleId)
		router.POST("/create", v1.CreateApi)
		router.PATCH("/:apiId", v1.UpdateApiById)
		router.DELETE("/batch", v1.BatchDeleteApiByIds)
	}
	return router
}
