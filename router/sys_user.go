package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 用户路由
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.POST("/info", v1.GetUserInfo)
		router.GET("/list", v1.GetUsers)
		router.PUT("/changePwd", v1.ChangePwd)
		router.POST("/create", v1.CreateUser)
		router.PATCH("/update/:userId", v1.UpdateUserById)
		router.DELETE("/delete/batch", v1.BatchDeleteUserByIds)
	}
	return router
}
