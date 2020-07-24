package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 文件上传路由
func InitUploadRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("upload").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/file", v1.UploadFileChunkExists)
		router.POST("/file", v1.UploadFile)
		router.POST("/merge", v1.UploadMerge)
		router.POST("/unzip", v1.UploadUnZip)
	}
	return router
}
