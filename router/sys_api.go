package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitApiRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/api")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/api")
	{
		router1.GET("/list", v1.GetApis)
		router1.GET("/all/category/:roleId", v1.GetAllApiGroupByCategoryByRoleId)
		router2.POST("/create", v1.CreateApi)
		router1.PATCH("/update/:apiId", v1.UpdateApiById)
		router1.DELETE("/delete/batch", v1.BatchDeleteApiByIds)
	}
	return r
}
