package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "go-shipment-api/api/v1"
	"go-shipment-api/middleware"
)

// 接口路由
func InitApiRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("api").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", v1.GetApis)
		router.POST("/create", v1.CreateApi)
		router.PATCH("/:apiId", v1.UpdateApiById)
		router.DELETE("/batch", v1.BatchDeleteApiByIds)
	}
	return router
}
