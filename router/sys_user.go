package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "go-shipment-api/api/v1"
)

// 用户路由
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("user").Use(authMiddleware.MiddlewareFunc())
	{
		router.POST("/info", v1.GetUserInfo)
		router.GET("/list", v1.GetUsers)
		router.PUT("/changePwd", v1.ChangePwd)
		router.POST("/create", v1.CreateUser)
		router.PATCH("/:userId", v1.UpdateUserById)
		router.DELETE("/batch", v1.BatchDeleteUserByIds)
	}
	return router
}
