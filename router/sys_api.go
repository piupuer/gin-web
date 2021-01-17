package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 接口路由
func InitApiRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/api")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/api")
	{
		router1.GET("/list", v1.GetApis)
		router1.GET("/all/category/:roleId", v1.GetAllApiGroupByCategoryByRoleId)
		router2.POST("/create", v1.CreateApi)
		router1.PATCH("/update/:apiId", v1.UpdateApiById)
		router1.DELETE("/delete/batch", v1.BatchDeleteApiByIds)
	}
	return r
}
