package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitUploadRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router := GetCasbinRouter(r, jwtOptions, "/upload")
	{
		router.GET("/file", v1.UploadFileChunkExists)
		router.POST("/file", v1.UploadFile)
		router.POST("/merge", v1.UploadMerge)
		router.POST("/unzip", v1.UploadUnZip)
	}
	return r
}
