package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 角色路由
func InitRoleRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/role")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/role")
	{
		router1.GET("/list", v1.GetRoles)
		router2.POST("/create", v1.CreateRole)
		router1.PATCH("/update/:roleId", v1.UpdateRoleById)
		router1.PATCH("/menus/update/:roleId", v1.UpdateRoleMenusById)
		router1.PATCH("/apis/update/:roleId", v1.UpdateRoleApisById)
		router1.DELETE("/delete/batch", v1.BatchDeleteRoleByIds)
	}
	return r
}
