package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 基础路由
func InitBaseRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("base")
	{
		router.POST("/login", authMiddleware.LoginHandler)
		router.POST("/refresh_token", authMiddleware.RefreshHandler)
	}
	return router
}
