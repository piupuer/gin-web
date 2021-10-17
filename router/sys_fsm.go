package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/middleware"
)

func InitFsmRouter(r *gin.RouterGroup, jwtOptions []func(*middleware.JwtOptions)) (R gin.IRoutes) {
	router1 := GetCasbinRouter(r, jwtOptions, "/fsm")
	router2 := GetCasbinAndIdempotenceRouter(r, jwtOptions, "/fsm")
	{
		router1.GET("/list", v1.FindFsm)
		router2.POST("/create", v1.CreateFsm)
		router1.GET("/approving/list", v1.FindFsmApprovingLog)
	}
	return r
}
