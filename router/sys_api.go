package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 接口路由
func InitApiRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("api").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", v1.GetApis)
		router.GET("/all/category/:roleId", v1.GetAllApiGroupByCategoryByRoleId)
		router.POST("/create", v1.CreateApi)
		router.PATCH("/update/:apiId", v1.UpdateApiById)
		router.DELETE("/delete/batch", v1.BatchDeleteApiByIds)
	}
	return router
}
