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
		router1.GET("/list", v1.FindDict)
		router2.POST("/create", v1.CreateDict)
		router1.PATCH("/update/:id", v1.UpdateDictById)
		router1.DELETE("/delete/batch", v1.BatchDeleteDictByIds)
		router1.GET("/data/list", v1.FindDictData)
		router2.POST("/data/create", v1.CreateDictData)
		router1.PATCH("/data/update/:id", v1.UpdateDictDataById)
		router1.DELETE("/data/delete/batch", v1.BatchDeleteDictDataByIds)
	}
	return r
}
