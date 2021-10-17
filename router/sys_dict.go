package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitDictRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/dict")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/dict")
	{
		router1.GET("/list", v1.GetDicts)
		router2.POST("/create", v1.CreateDict)
		router1.PATCH("/update/:dictId", v1.UpdateDictById)
		router1.DELETE("/delete/batch", v1.BatchDeleteDictByIds)
		router1.GET("/data/list", v1.GetDictDatas)
		router2.POST("/data/create", v1.CreateDictData)
		router1.PATCH("/data/update/:dictDataId", v1.UpdateDictDataById)
		router1.DELETE("/data/delete/batch", v1.BatchDeleteDictDataByIds)
	}
	return r
}
