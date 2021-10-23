package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitRoleRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/role")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/role")
	{
		router1.GET("/list", v1.FindRole)
		router2.POST("/create", v1.CreateRole)
		router1.PATCH("/update/:id", v1.UpdateRoleById)
		router1.DELETE("/delete/batch", v1.BatchDeleteRoleByIds)
	}
	return r
}
