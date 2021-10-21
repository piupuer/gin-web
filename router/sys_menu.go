package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitMenuRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/menu")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/menu")
	{
		router1.GET("/tree", v1.GetMenuTree)
		router1.GET("/all/:id", v1.GetAllMenuByRoleId)
		router1.GET("/list", v1.FindMenu)
		router2.POST("/create", v1.CreateMenu)
		router1.PATCH("/update/:id", v1.UpdateMenuById)
		router2.DELETE("/delete/batch", v1.BatchDeleteMenuByIds)
	}
	return r
}
