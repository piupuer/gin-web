package router

import (
	v1 "gin-web/api/v1"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 用户路由
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/user")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/user")
	{
		router1.POST("/info", v1.GetUserInfo)
		router1.GET("/list", v1.GetUsers)
		router1.PUT("/changePwd", v1.ChangePwd)
		router2.POST("/create", v1.CreateUser)
		router1.PATCH("/update/:userId", v1.UpdateUserById)
		router1.DELETE("/delete/batch", v1.BatchDeleteUserByIds)
	}
	return r
}
