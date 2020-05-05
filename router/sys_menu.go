package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "go-shipment-api/api/v1"
	"go-shipment-api/middleware"
)

// 菜单路由
func InitMenuRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("menu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/tree", v1.GetMenuTree)
		router.GET("/list", v1.GetMenus)
		router.POST("/create", v1.CreateMenu)
		router.PATCH("/:menuId", v1.UpdateMenuById)
		router.DELETE("/batch", v1.BatchDeleteMenuByIds)
	}
	return router
}
