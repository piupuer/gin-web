package router

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 字典路由
func InitDictRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/dict").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", v1.GetDicts)
		router.POST("/create", v1.CreateDict)
		router.PATCH("/update/:dictId", v1.UpdateDictById)
		router.DELETE("/delete/batch", v1.BatchDeleteDictByIds)
		router.GET("/data/list", v1.GetDictDatas)
		router.POST("/data/create", v1.CreateDictData)
		router.PATCH("/data/update/:dictDataId", v1.UpdateDictDataById)
		router.DELETE("/data/delete/batch", v1.BatchDeleteDictDataByIds)
	}
	return router
}
