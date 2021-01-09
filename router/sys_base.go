package router

import (
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 基础路由
func InitBaseRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/base")
	{
		// 登录登出/刷新token无需鉴权
		router.POST("/login", authMiddleware.LoginHandler)
		router.POST("/logout", authMiddleware.LogoutHandler)
		router.POST("/refreshToken", authMiddleware.RefreshHandler)
		// 幂等性token需要鉴权
		router.
			Use(authMiddleware.MiddlewareFunc()).
			Use(middleware.CasbinMiddleware).
			GET("/idempotenceToken", middleware.GetIdempotenceToken)
	}
	return r
}

// 获取带casbin中间件的路由
func GetCasbinRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware, path string) gin.IRoutes {
	return r.Group(path).Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
}

// 获取带casbin和幂等性中间件的路由
func GetCasbinAndIdempotenceRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware, path string) gin.IRoutes {
	return GetCasbinRouter(r, authMiddleware, path).Use(middleware.Idempotence)
}
